package service
import (
 "testing"
 "github.com/contractapi/contractapi/internal/model"
 "github.com/contractapi/contractapi/internal/repository"
)
func TestTicketReplyChain(t *testing.T){
 db:=newCtDB(t)
 tr:=repository.NewTicketRepository(db)
 svc:=NewTicketService(tr,ctLogger())
 tk:=&model.LegalTicket{UserID:1,Type:"consult",Title:"t",Status:"pending"}; tr.Create(tk)
 if _,err:=svc.AddReply(1,tk.ID,"user","hi",nil);err!=nil{t.Fatalf("AddReply: %v",err)}
}
