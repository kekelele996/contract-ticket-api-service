package service
import (
 "testing"
 "github.com/contractapi/contractapi/internal/repository"
)
func expectNoPanicCt(t *testing.T,name string,fn func() error) error { t.Helper(); var err error; func(){ defer func(){ if r:=recover(); r!=nil { t.Fatalf("%s panicked: %v",name,r) } }(); err=fn() }(); return err }
func TestMissingContractAndTicket(t *testing.T){
 db:=newCtDB(t)
 cr:=repository.NewContractRepository(db); tr:=repository.NewTicketRepository(db); tm:=repository.NewTemplateRepository(db)
 cSvc:=NewContractService(cr,tm,NewPDFService(ctLogger()),ctLogger())
 tSvc:=NewTicketService(tr,ctLogger())
 if err:=expectNoPanicCt(t,"Submit",func() error { return cSvc.Submit(1,999,nil) }); err==nil{t.Fatal("Submit should return error")}
 if err:=expectNoPanicCt(t,"AddReply",func() error { _,err:=tSvc.AddReply(1,999,"user","hi",nil); return err }); err==nil{t.Fatal("AddReply should return error")}
}
