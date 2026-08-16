package service
import (
 "testing"
 "github.com/contractapi/contractapi/internal/dto"
 "github.com/contractapi/contractapi/internal/model"
 "github.com/contractapi/contractapi/internal/repository"
)
func TestContractSubmitAndSign(t *testing.T){
 db:=newCtDB(t)
 cr:=repository.NewContractRepository(db); tr:=repository.NewTemplateRepository(db)
 svc:=NewContractService(cr,tr,NewPDFService(ctLogger()),ctLogger())
 c:=&model.Contract{UserID:1,TemplateID:1,Title:"合同",Status:"draft",Variables:model.JSONMap{}}
 if err:=cr.Create(c);err!=nil{t.Fatal(err)}
 if err:=svc.Submit(1,c.ID,[]dto.SignerInput{{Name:"甲",Role:"party_a"}});err!=nil{t.Fatalf("Submit: %v",err)}
 if err:=svc.Sign(1,c.ID,"甲","party_a","");err!=nil{t.Fatalf("Sign: %v",err)}
}
