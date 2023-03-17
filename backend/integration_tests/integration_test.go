package integration_tests

import (
	"bytes"
	"data-graph-backend/pkg/apiServer"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/logging"
	"data-graph-backend/pkg/properties"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	_, b, _, _ = runtime.Caller(0)
	configPath = filepath.Dir(b) + "/../config/config_test.yaml"
	router     = InitializeConnection(configPath)
)

func InitializeConnection(configPath string) *apiServer.Router {
	config, err := properties.GetConfig(configPath)
	if err != nil {
		log.Fatal("Can't get config at the path: " + configPath)
	}

	logging.Init(config.ProgramSettings.LogPath)
	logger := logging.GetLogger()
	logger.Info("Starting backend for DataGraph")

	dbConnection, err := dbConnector.NewConnection(config, logger)
	if err != nil {
		log.Fatal("Can't connect to db - ", err.Error())
	}
	err = dbConnection.SetIdShift()
	if err != nil {
		log.Fatal("Can't get info from db - ", err.Error())
	}

	router := &apiServer.Router{
		Logger:      logger,
		DbConnector: dbConnection,
	}

	apiServer.ConfigureRouters(router)
	return router
}

func Test_successful_get_company_data(t *testing.T) {
	expected := "{\"id\":80,\"name\":\"Яндекс\",\"ceo\":\"Артём Савиновский\",\"description\":\"«Яндекс» - транснациональная компания в отрасли информационных технологий, чьё головное юридическое лицо зарегистрировано в Нидерландах, владеющая одноимённой системой поиска в интернете, интернет-порталом и веб-службами в нескольких странах. Наиболее заметное положение занимает на рынках России, Белоруссии и Казахстана\",\"staffSize\":10227,\"year\":\"2000-01-01T00:00:00Z\",\"departments\":[\"Транспорт\",\"Связь\",\"Финансы\",\"IT\",\"Общественное питание\",\"Торговля\",\"Развлечения\"],\"products\":[{\"id\":2,\"name\":\"Яndex-Web\",\"year\":\"1997-09-23T00:00:00Z\",\"isVerified\":false},{\"id\":4,\"name\":\"Яndex-Web\",\"year\":\"1998-10-16T00:00:00Z\",\"isVerified\":false},{\"id\":6,\"name\":\"Яndex-Web\",\"year\":\"1999-06-09T00:00:00Z\",\"isVerified\":false},{\"id\":8,\"name\":\"Яndex-Web\",\"year\":\"2001-06-13T00:00:00Z\",\"isVerified\":false},{\"id\":10,\"name\":\"Яndex-Web\",\"year\":\"2005-12-23T00:00:00Z\",\"isVerified\":false},{\"id\":12,\"name\":\"Яndex-Web\",\"year\":\"2010-04-22T00:00:00Z\",\"isVerified\":false},{\"id\":14,\"name\":\"Яndex-Web\",\"year\":\"2011-08-19T00:00:00Z\",\"isVerified\":false},{\"id\":20,\"name\":\"Яndex-Web\",\"year\":\"2014-06-13T00:00:00Z\",\"isVerified\":false},{\"id\":16,\"name\":\"Яндекс Браузер\",\"year\":\"2012-10-01T00:00:00Z\",\"isVerified\":false},{\"id\":18,\"name\":\"Яндекс Браузер\",\"year\":\"2015-06-18T00:00:00Z\",\"isVerified\":false},{\"id\":22,\"name\":\"Яндекс Дзен\",\"year\":\"2015-06-10T00:00:00Z\",\"isVerified\":false},{\"id\":24,\"name\":\"Яндекс Дзен\",\"year\":\"2018-08-16T00:00:00Z\",\"isVerified\":false},{\"id\":26,\"name\":\"Яндекс Дзен\",\"year\":\"2021-10-13T00:00:00Z\",\"isVerified\":false}],\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\"}"

	svr := httptest.NewServer(http.HandlerFunc(router.HandleCompany))
	defer svr.Close()

	resp, err := http.Get(svr.URL + "?id=80")
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_successful_get_product_data(t *testing.T) {
	expected := "{\"id\":4,\"name\":\"Яndex-Web\",\"version\":\"1.1\",\"company\":{\"id\":3,\"name\":\"Яндекс\"},\"link\":\"\",\"description\":\"Яндекс Поисковик\",\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"year\":\"1998-10-16T00:00:00Z\",\"departments\":[{\"id\":123,\"name\":\"IT\"}]}"

	svr := httptest.NewServer(http.HandlerFunc(router.HandleProduct))
	defer svr.Close()

	resp, err := http.Get(svr.URL + "?id=4")
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_successful_get_filter_company_data(t *testing.T) {
	expected := "[79,28,30,32,34,36,38,40]"

	svr := httptest.NewServer(http.HandlerFunc(router.HandleFilterCompany))
	defer svr.Close()

	filtersJson := []byte("{\"companyName\":\"VK\",\"departments\":[1,2,3],\"ceo\":\"\",\"minDate\":\"1996-01-02T15:04:05Z\",\"maxDate\":\"2020-01-02T15:04:05Z\",\"startStaffSize\":100,\"endStaffSize\":500000}")
	resp, err := http.Post(svr.URL, "application/json", bytes.NewBuffer(filtersJson))
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_successful_get_graph_data(t *testing.T) {
	expected := "{\"nodes\":[{\"name\":\"Яндекс\",\"id\":80,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Компания\",\"size\":900,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":0},{\"name\":\"VK\",\"id\":79,\"svg\":\"http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png\",\"nodeType\":\"Компания\",\"size\":900,\"opacity\":1,\"color\":\"#8AC926\",\"x\":100,\"y\":0},{\"name\":\"Ozon\",\"id\":81,\"svg\":\"http://141.95.127.215:8082/media-server/data/3a39568e528f3a06794cffde50e10730.png\",\"nodeType\":\"Компания\",\"size\":900,\"opacity\":1,\"color\":\"#6A4C93\",\"x\":1100,\"y\":0},{\"name\":\"Фирма 1С\",\"id\":82,\"svg\":\"http://141.95.127.215:8082/media-server/data/768c058be0c89e5f563250e1134a0ab1.png\",\"nodeType\":\"Компания\",\"size\":900,\"opacity\":1,\"color\":\"#FF595E\",\"x\":1600,\"y\":0},{\"name\":\"СБЕР\",\"id\":78,\"svg\":\"http://141.95.127.215:8082/media-server/data/20ad243690efeb01588568f452f8593a.jpg\",\"nodeType\":\"Компания\",\"size\":900,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-400,\"y\":0},{\"name\":\"Яndex-Web\",\"id\":2,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":200},{\"name\":\"Яndex-Web\",\"id\":4,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":300},{\"name\":\"Яndex-Web\",\"id\":6,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":400},{\"name\":\"Яndex-Web\",\"id\":8,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":500},{\"name\":\"Яndex-Web\",\"id\":10,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":600},{\"name\":\"Яndex-Web\",\"id\":12,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":700},{\"name\":\"Яndex-Web\",\"id\":14,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":800},{\"name\":\"Яndex-Web\",\"id\":20,\"svg\":\"http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":600,\"y\":1100},{\"name\":\"Яндекс Браузер\",\"id\":16,\"svg\":\"http://141.95.127.215:8082/media-server/data/65c33c46d41707be281b1d93115c2136.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":700,\"y\":900},{\"name\":\"Яндекс Браузер\",\"id\":18,\"svg\":\"http://141.95.127.215:8082/media-server/data/65c33c46d41707be281b1d93115c2136.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":700,\"y\":1000},{\"name\":\"Яндекс Дзен\",\"id\":22,\"svg\":\"http://141.95.127.215:8082/media-server/data/57968135f081664fc2c0fc04a643a713.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":500,\"y\":1200},{\"name\":\"Яндекс Дзен\",\"id\":24,\"svg\":\"http://141.95.127.215:8082/media-server/data/57968135f081664fc2c0fc04a643a713.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":500,\"y\":1300},{\"name\":\"Яндекс Дзен\",\"id\":26,\"svg\":\"http://141.95.127.215:8082/media-server/data/57968135f081664fc2c0fc04a643a713.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#1982C4\",\"x\":500,\"y\":1400},{\"name\":\"Дзен\",\"id\":28,\"svg\":\"http://141.95.127.215:8082/media-server/data/57968135f081664fc2c0fc04a643a713.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":300,\"y\":1500},{\"name\":\"Яндекс Дзен\",\"id\":30,\"svg\":\"http://141.95.127.215:8082/media-server/data/57968135f081664fc2c0fc04a643a713.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":300,\"y\":1600},{\"name\":\"Vkontakte\",\"id\":32,\"svg\":\"http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":-50,\"y\":200},{\"name\":\"Vkontakte\",\"id\":34,\"svg\":\"http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":-50,\"y\":300},{\"name\":\"Vkontakte\",\"id\":36,\"svg\":\"http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":-50,\"y\":400},{\"name\":\"Mail.ru\",\"id\":38,\"svg\":\"http://141.95.127.215:8082/media-server/data/444ae1da2700275bdcc9a2ae397a9cfc.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":350,\"y\":200},{\"name\":\"Mail.ru\",\"id\":40,\"svg\":\"http://141.95.127.215:8082/media-server/data/444ae1da2700275bdcc9a2ae397a9cfc.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#8AC926\",\"x\":350,\"y\":300},{\"name\":\"Ozon.ru\",\"id\":42,\"svg\":\"http://141.95.127.215:8082/media-server/data/3a39568e528f3a06794cffde50e10730.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#6A4C93\",\"x\":900,\"y\":100},{\"name\":\"Ozon.ru\",\"id\":44,\"svg\":\"http://141.95.127.215:8082/media-server/data/3a39568e528f3a06794cffde50e10730.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#6A4C93\",\"x\":900,\"y\":200},{\"name\":\"Ozon.ru\",\"id\":67,\"svg\":\"http://141.95.127.215:8082/media-server/data/3a39568e528f3a06794cffde50e10730.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#6A4C93\",\"x\":900,\"y\":300},{\"name\":\"OzonApp\",\"id\":48,\"svg\":\"http://141.95.127.215:8082/media-server/data/3a39568e528f3a06794cffde50e10730.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#6A4C93\",\"x\":1000,\"y\":300},{\"name\":\"OzonApp\",\"id\":50,\"svg\":\"http://141.95.127.215:8082/media-server/data/3a39568e528f3a06794cffde50e10730.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#6A4C93\",\"x\":1000,\"y\":400},{\"name\":\"SberbankOnline\",\"id\":52,\"svg\":\"http://141.95.127.215:8082/media-server/data/2d49f9618e2c3252d205df0ea717088d.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-500,\"y\":200},{\"name\":\"SberbankOnline\",\"id\":54,\"svg\":\"http://141.95.127.215:8082/media-server/data/2d49f9618e2c3252d205df0ea717088d.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-500,\"y\":300},{\"name\":\"SberbankOnline\",\"id\":56,\"svg\":\"http://141.95.127.215:8082/media-server/data/2d49f9618e2c3252d205df0ea717088d.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-500,\"y\":400},{\"name\":\"SberApp\",\"id\":59,\"svg\":\"http://141.95.127.215:8082/media-server/data/20ad243690efeb01588568f452f8593a.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-300,\"y\":200},{\"name\":\"SberApp\",\"id\":61,\"svg\":\"http://141.95.127.215:8082/media-server/data/20ad243690efeb01588568f452f8593a.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-300,\"y\":300},{\"name\":\"SberApp\",\"id\":65,\"svg\":\"http://141.95.127.215:8082/media-server/data/20ad243690efeb01588568f452f8593a.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FFCA3A\",\"x\":-300,\"y\":400},{\"name\":\"1С.Предприятия\",\"id\":71,\"svg\":\"http://141.95.127.215:8082/media-server/data/d721dfb68da72edccba5e902687c2702.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FF595E\",\"x\":1600,\"y\":200},{\"name\":\"1С.Бухгалтерия\",\"id\":73,\"svg\":\"http://141.95.127.215:8082/media-server/data/945ff2dbf16a47e7be5a48e4005de553.jpg\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FF595E\",\"x\":1600,\"y\":400},{\"name\":\"1С.Битрикс\",\"id\":75,\"svg\":\"http://141.95.127.215:8082/media-server/data/7649495a9176b5aa3b98d616f6c6aa26.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FF595E\",\"x\":1400,\"y\":300},{\"name\":\"1С.Розница\",\"id\":77,\"svg\":\"http://141.95.127.215:8082/media-server/data/818f5d89989e195d35d3a6a78df31a11.png\",\"nodeType\":\"Продукт\",\"size\":600,\"opacity\":1,\"color\":\"#FF595E\",\"x\":1800,\"y\":300}],\"links\":[{\"source\":80,\"target\":2,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":2,\"target\":4,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":4,\"target\":6,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":6,\"target\":8,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":8,\"target\":10,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":10,\"target\":12,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":12,\"target\":14,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":14,\"target\":20,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":14,\"target\":16,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":16,\"target\":18,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":20,\"target\":22,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":22,\"target\":24,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":24,\"target\":26,\"color\":\"#1982C4\",\"opacity\":1},{\"source\":26,\"target\":28,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":79,\"target\":28,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":28,\"target\":30,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":79,\"target\":32,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":32,\"target\":34,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":34,\"target\":36,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":79,\"target\":38,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":38,\"target\":40,\"color\":\"#8AC926\",\"opacity\":1},{\"source\":81,\"target\":42,\"color\":\"#6A4C93\",\"opacity\":1},{\"source\":42,\"target\":44,\"color\":\"#6A4C93\",\"opacity\":1},{\"source\":44,\"target\":67,\"color\":\"#6A4C93\",\"opacity\":1},{\"source\":44,\"target\":48,\"color\":\"#6A4C93\",\"opacity\":1},{\"source\":48,\"target\":50,\"color\":\"#6A4C93\",\"opacity\":1},{\"source\":78,\"target\":52,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":52,\"target\":54,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":54,\"target\":56,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":78,\"target\":59,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":54,\"target\":61,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":59,\"target\":61,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":61,\"target\":65,\"color\":\"#FFCA3A\",\"opacity\":1},{\"source\":82,\"target\":71,\"color\":\"#FF595E\",\"opacity\":1},{\"source\":71,\"target\":73,\"color\":\"#FF595E\",\"opacity\":1},{\"source\":71,\"target\":75,\"color\":\"#FF595E\",\"opacity\":1},{\"source\":71,\"target\":77,\"color\":\"#FF595E\",\"opacity\":1}]}"
	svr := httptest.NewServer(http.HandlerFunc(router.HandleGetGraphFull))

	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_successful_get_filter_product_data(t *testing.T) {
	expected := "[32,34,36,38,40,52,54,56,59,61,71,73,75,77]"

	svr := httptest.NewServer(http.HandlerFunc(router.HandleFilterProduct))
	defer svr.Close()

	filtersJson := []byte("{\"productName\":\"\",\"minDate\":\"1996-01-02T15:04:05Z\",\"maxDate\":\"2020-01-02T15:04:05Z\",\"departments\":[1,2,3],\"isVerified\":false}")
	resp, err := http.Post(svr.URL, "application/json", bytes.NewBuffer(filtersJson))
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_successful_get_company_filter_preset(t *testing.T) {
	expected := "{\"companyFilters\":{\"companyNames\":[\"Яндекс\",\"VK\",\"Ozon\",\"СБЕР\",\"Фирма 1С\"],\"ceoNames\":[\"Герман Оскарович Греф\",\"Нуралиев Борис Георгиевич\",\"Артём Савиновский\",\"Сергей Паньков\",\"Владимир Сергеевич Кириенко\"],\"minStaffSize\":1100,\"maxStaffSize\":285000,\"minDate\":\"1991-01-01T00:00:00Z\",\"maxDate\":\"2000-01-01T00:00:00Z\",\"departments\":[{\"id\":1,\"name\":\"Транспорт\"},{\"id\":2,\"name\":\"Связь\"},{\"id\":3,\"name\":\"Финансы\"},{\"id\":4,\"name\":\"IT\"},{\"id\":5,\"name\":\"Общественное питание\"},{\"id\":6,\"name\":\"Торговля\"},{\"id\":7,\"name\":\"Фармацевтика\"},{\"id\":8,\"name\":\"Развлечения\"}]},\"productFilters\":{\"productNames\":[\"SberbankOnline\",\"Ozon.ru\",\"1С.Бухгалтерия\",\"OzonApp\",\"Mail.ru\",\"Яндекс Браузер\",\"SberApp\",\"1С.Предприятия\",\"Vkontakte\",\"Яndex-Web\",\"Яндекс Дзен\",\"1С.Розница\",\"1С.Битрикс\",\"Дзен\"],\"minDate\":\"1996-12-13T00:00:00Z\",\"maxDate\":\"2022-12-20T00:00:00Z\"}}"

	svr := httptest.NewServer(http.HandlerFunc(router.HandleGetFilterPresets))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_successful_get_filter_preset(t *testing.T) {
	expected := "{\"companyFilters\":{\"companyNames\":[\"Яндекс\",\"VK\",\"Ozon\",\"СБЕР\",\"Фирма 1С\"],\"ceoNames\":[\"Герман Оскарович Греф\",\"Нуралиев Борис Георгиевич\",\"Артём Савиновский\",\"Сергей Паньков\",\"Владимир Сергеевич Кириенко\"],\"minStaffSize\":1100,\"maxStaffSize\":285000,\"minDate\":\"1991-01-01T00:00:00Z\",\"maxDate\":\"2000-01-01T00:00:00Z\",\"departments\":[{\"id\":1,\"name\":\"Транспорт\"},{\"id\":2,\"name\":\"Связь\"},{\"id\":3,\"name\":\"Финансы\"},{\"id\":4,\"name\":\"IT\"},{\"id\":5,\"name\":\"Общественное питание\"},{\"id\":6,\"name\":\"Торговля\"},{\"id\":7,\"name\":\"Фармацевтика\"},{\"id\":8,\"name\":\"Развлечения\"}]},\"productFilters\":{\"productNames\":[\"SberbankOnline\",\"Ozon.ru\",\"1С.Бухгалтерия\",\"OzonApp\",\"Mail.ru\",\"Яндекс Браузер\",\"SberApp\",\"1С.Предприятия\",\"Vkontakte\",\"Яndex-Web\",\"Яндекс Дзен\",\"1С.Розница\",\"1С.Битрикс\",\"Дзен\"],\"minDate\":\"1996-12-13T00:00:00Z\",\"maxDate\":\"2022-12-20T00:00:00Z\"}}"

	svr := httptest.NewServer(http.HandlerFunc(router.HandleGetFilterPresets))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	res = strings.TrimSpace(res)
	res = strings.TrimLeft(res, "\"")
	res = strings.TrimRight(res, "\"")
	if res != expected {
		t.Errorf("expected res to be %s \ngot %s", expected, res)
	}

}

func Test_unsuccessful_get_company_by_id(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(router.HandleCompany))
	defer svr.Close()

	resp, err := http.Get(svr.URL + "?id=8000")
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	if resp.Status != "200 OK" {
		if !strings.Contains(res, "Wrong argument") {
			t.Errorf("Unknown error: %s", err.Error())
		}
		return
	}
	t.Errorf("Don't get an error")
}

func Test_unsuccessful_get_product_by_id(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(router.HandleProduct))
	defer svr.Close()

	resp, err := http.Get(svr.URL + "?id=8000")
	if err != nil {
		log.Fatal(err, "unable to complete Get request")
	}
	resB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "unable to read response data")
	}
	res := string(resB)

	if resp.Status != "200 OK" {
		if !strings.Contains(res, "Wrong argument") {
			t.Errorf("Unknown error: %s", err.Error())
		}
		return
	}
	t.Errorf("Don't get an error")
}
