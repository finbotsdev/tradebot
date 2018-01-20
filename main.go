package main

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/strategy"
	"github.com/jeremyhahn/tradebot/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

func main() {

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	sqlite := InitSQLite()
	defer sqlite.Close()

	ctx := &common.Context{
		DB:     sqlite,
		Logger: logger}

	userDAO := dao.NewUserDAO(ctx)
	ctx.User = userDAO.GetById(1)

	marketcapService := service.NewMarketCapService(logger)

	ws := websocket.NewWebsocketServer(ctx, 8080, marketcapService)
	go ws.Start()

	//tradeService := service.NewTradeService(ctx, marketcapService)
	//tradeService.Trade()

	var services []common.ChartService
	exchangeDAO := dao.NewExchangeDAO(ctx)
	autoTradeDAO := dao.NewAutoTradeDAO(ctx)
	signalDAO := dao.NewSignalLogDAO(ctx)
	for _, autoTradeCoin := range autoTradeDAO.Find(ctx.User) {
		ctx.Logger.Debugf("[NewTradeService] Loading AutoTrade currency pair: %s-%s\n", autoTradeCoin.Base, autoTradeCoin.Quote)
		currencyPair := &common.CurrencyPair{
			Base:          autoTradeCoin.Base,
			Quote:         autoTradeCoin.Quote,
			LocalCurrency: ctx.User.LocalCurrency}
		exchangeService := service.NewExchangeService(ctx, exchangeDAO)
		exchange := exchangeService.NewExchange(ctx.User, autoTradeCoin.Exchange, currencyPair)
		strategy := strategy.NewDefaultTradingStrategy(ctx, &autoTradeCoin, autoTradeDAO, signalDAO)
		chart := service.NewChartService(ctx, exchange, strategy, autoTradeCoin.Period)
		ctx.Logger.Debugf("[NewTradeService] Chart: %+v\n", chart)
		services = append(services, chart)
	}

	for _, chart := range services {
		chart.Stream()
	}

}

func InitSQLite() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/tradebot.db")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}

/*
func InitMySQL() *gorm.DB {
	db, err := gorm.Open("mysql", "user:pass@tcp(ip:3306)/mydb?charset=utf8&parseTime=True")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}
*/
