package reporters

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"

	"github.com/onflow/cadence/common"
	"github.com/onflow/cadence/interpreter"
	cadenceRuntime "github.com/onflow/cadence/runtime"

	"github.com/onflow/flow-go/cmd/util/ledger/util"
	"github.com/onflow/flow-go/fvm/environment"
	"github.com/onflow/flow-go/fvm/storage/state"
	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/ledger"
	"github.com/onflow/flow-go/model/flow"
)

const FungibleTokenTrackerReportPrefix = "fungible_token_report"

var domains = []common.StorageDomain{
	common.StorageDomainPathPublic,
	common.StorageDomainPathPrivate,
	common.StorageDomainPathStorage,
}

// FungibleTokenTracker iterates through stored cadence values over all accounts and check for any
// value with the given resource typeID
type FungibleTokenTracker struct {
	log          zerolog.Logger
	chain        flow.Chain
	rwf          ReportWriterFactory
	rw           ReportWriter
	progress     *progressbar.ProgressBar
	vaultTypeIDs map[string]bool
}

func FlowTokenTypeID(chain flow.Chain) string {
	sc := systemcontracts.SystemContractsForChain(chain.ChainID())
	return fmt.Sprintf("A.%s.FlowToken.Vault", sc.FlowToken.Address.Hex())
}

func NewFungibleTokenTracker(logger zerolog.Logger, rwf ReportWriterFactory, chain flow.Chain, vaultTypeIDs []string) *FungibleTokenTracker {
	ftt := &FungibleTokenTracker{
		log:          logger,
		rwf:          rwf,
		chain:        chain,
		vaultTypeIDs: make(map[string]bool),
	}
	for _, vt := range vaultTypeIDs {
		ftt.vaultTypeIDs[vt] = true
	}
	return ftt
}

func (r *FungibleTokenTracker) Name() string {
	return "Resource Tracker"
}

type trace []string

func (t trace) String() string {
	return strings.Join(t, "/")
}

type TokenDataPoint struct {
	// Path is the storage path of the composite the vault was found in
	Path string `json:"path"`
	// Address is the owner of the composite the vault was found in
	Address string `json:"address"`
	// Balance is the balance of the flow vault
	Balance uint64 `json:"balance"`
	// token type
	TypeID string `json:"type_id"`
}

type job struct {
	owner    flow.Address
	payloads []ledger.Payload
}

// Report creates a fungible_token_report_*.json file that contains data on all fungible token Vaults in the state commitment.
// I recommend using gojq to browse through the data, because of the large uint64 numbers which jq won't be able to handle.
func (r *FungibleTokenTracker) Report(payloads []ledger.Payload, commit ledger.State) error {
	r.rw = r.rwf.ReportWriter(FungibleTokenTrackerReportPrefix)
	defer r.rw.Close()

	wg := &sync.WaitGroup{}

	// we need to shard by owner, otherwise ledger won't be thread-safe
	addressCount := 0
	payloadsByOwner := make(map[flow.Address][]ledger.Payload)

	for _, pay := range payloads {
		k, err := pay.Key()
		if err != nil {
			return nil
		}
		owner := flow.BytesToAddress(k.KeyParts[0].Value)
		if len(owner) > 0 { // ignoring payloads without ownership (fvm ones)
			m, ok := payloadsByOwner[owner]
			if !ok {
				payloadsByOwner[owner] = make([]ledger.Payload, 0)
				addressCount++
			}
			payloadsByOwner[owner] = append(m, pay)
		}
	}

	jobs := make(chan job, addressCount)
	r.progress = progressbar.Default(int64(addressCount), "Processing:")

	for k, v := range payloadsByOwner {
		jobs <- job{k, v}
	}

	close(jobs)

	workerCount := runtime.NumCPU()
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go r.worker(jobs, wg)
	}

	wg.Wait()

	err := r.progress.Finish()
	if err != nil {
		panic(err)
	}

	return nil
}

func (r *FungibleTokenTracker) worker(
	jobs <-chan job,
	wg *sync.WaitGroup) {
	for j := range jobs {

		inter, err := interpreter.NewInterpreter(nil, nil, &interpreter.Config{})
		if err != nil {
			panic(err)
		}

		txnState := state.NewTransactionState(
			NewStorageSnapshotFromPayload(j.payloads),
			state.DefaultParameters())
		accounts := environment.NewAccounts(txnState)
		storage := cadenceRuntime.NewStorage(
			&util.AccountsAtreeLedger{Accounts: accounts},
			nil,
			cadenceRuntime.StorageConfig{},
		)

		owner, err := common.BytesToAddress(j.owner[:])
		if err != nil {
			panic(err)
		}

		for _, domain := range domains {
			storageMap := storage.GetDomainStorageMap(inter, owner, domain, true)
			itr := storageMap.Iterator(inter)
			key, value := itr.Next()
			for value != nil {
				identifier := string(key.(interpreter.StringAtreeValue))
				r.iterateChildren(
					append([]string{domain.Identifier()}, identifier),
					j.owner,
					value,
				)
				key, value = itr.Next()
			}
		}

		err = r.progress.Add(1)
		if err != nil {
			panic(err)
		}
	}

	wg.Done()
}

func (r *FungibleTokenTracker) iterateChildren(tr trace, addr flow.Address, value interpreter.Value) {

	compValue, ok := value.(*interpreter.CompositeValue)
	if !ok {
		return
	}

	// because compValue.Kind == common.CompositeKindResource
	// we could pass nil to the IsResourceKinded method
	inter, err := interpreter.NewInterpreter(nil, nil, &interpreter.Config{})
	if err != nil {
		panic(err)
	}
	if compValue.IsResourceKinded(nil) {
		typeIDStr := string(compValue.TypeID())
		if _, ok := r.vaultTypeIDs[typeIDStr]; ok {
			b := uint64(compValue.GetField(
				inter,
				"balance",
			).(interpreter.UFix64Value).UFix64Value)
			if b > 0 {
				r.rw.Write(TokenDataPoint{
					Path:    tr.String(),
					Address: addr.Hex(),
					Balance: b,
					TypeID:  string(compValue.TypeID()),
				})
			}
		}

		// iterate over fields of the composite value (skip the ones that are not resource typed)
		compValue.ForEachField(inter,
			func(key string, value interpreter.Value) (resume bool) {
				r.iterateChildren(append(tr, key), addr, value)

				// continue iteration
				return true
			},
			interpreter.EmptyLocationRange,
		)
	}
}
