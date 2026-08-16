package service
import (
 "errors"
 "testing"
 "github.com/contractapi/contractapi/internal/model"
 "github.com/contractapi/contractapi/internal/repository"
)
func TestContractListByUser(t *testing.T){
 db:=newCtDB(t)
 cr:=repository.NewContractRepository(db)
 cr.Create(&model.Contract{UserID:1,TemplateID:1,Title:"a",Status:"draft"})
 c2:=&model.Contract{UserID:2,TemplateID:1,Title:"b",Status:"draft"}; cr.Create(c2)
 if _,err:=cr.FindByIDForUser(c2.ID,1); !errors.Is(err,repository.ErrNotFound){t.Fatalf("other user contract should be not found, got %v",err)}
 list,total,_:=cr.ListByUser(1,"",0,10)
 if total!=1 || len(list)!=1 {t.Fatalf("list = %d/%d, want 1",len(list),total)}
}
