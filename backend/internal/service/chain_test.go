package service
import (
 "errors"
 "testing"
 "github.com/contractapi/contractapi/internal/repository"
)
func TestErrorChainSentinels(t *testing.T){
 db:=newCtDB(t)
 cr:=repository.NewContractRepository(db); tr:=repository.NewTicketRepository(db); tm:=repository.NewTemplateRepository(db)
 if _,err:=cr.FindByID(999);!errors.Is(err,repository.ErrNotFound){t.Fatalf("contract=%v",err)}
 if _,err:=tr.FindByID(999);!errors.Is(err,repository.ErrNotFound){t.Fatalf("ticket=%v",err)}
 if _,err:=tm.FindByID(999);!errors.Is(err,repository.ErrNotFound){t.Fatalf("template=%v",err)}
}
