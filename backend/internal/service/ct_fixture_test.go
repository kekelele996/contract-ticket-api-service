package service

import (
	"log/slog"
	"os"
	"testing"

	"github.com/contractapi/contractapi/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newCtDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil { t.Fatalf("open sqlite: %v", err) }
	if err := db.AutoMigrate(&model.User{}, &model.Contract{}, &model.ContractSigner{}, &model.ContractTemplate{}, &model.LegalTicket{}, &model.TicketReply{}, &model.KnowledgeFAQ{}, &model.TemplateFavorite{}); err != nil { t.Fatalf("migrate: %v", err) }
	return db
}

func ctLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}
