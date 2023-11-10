package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getAppJourney(c *gin.Context) {
	var af AppFilter

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		msg := `id invalid or missing`
		fmt.Println(msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	af.AppID = id

	if err := c.ShouldBindQuery(&af); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := af.validate(); err != nil {
		msg := "app journey request validation failed"
		fmt.Println(msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})
		return
	}

	if !af.hasTimeRange() {
		af.setDefaultTimeRange()
	}

	if !af.hasVersion() {
		af.setDefaultVersion()
	}

	// fmt.Println("journey request app id", af.AppID)
	// fmt.Println("journey request from", af.From)
	// fmt.Println("journey request to", af.To)
	// fmt.Println("journey request version", af.Version)

	data1 := `{"nodes":[{"id":"Home Screen","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Order History","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Order Status","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Support","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"List Of Items","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Sales Offer","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"View Item Images","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"View Item Detail","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Cyber Monday Sale Items List","nodeColor":"hsl(0, 72%, 51%)","issues":{"crashes":[{"title":"NullPointerException.java","count":37893},{"title":"LayoutInflaterException.java","count":12674}],"anrs":[{"title":"CyberMondayActivity.java","count":97321},{"title":"CyberMondayFragment.kt","count":8005}]}},{"id":"Add To Cart","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Pay","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Explore Discounts","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}}],"links":[{"source":"Home Screen","target":"Order History","value":50000},{"source":"Home Screen","target":"List Of Items","value":73356},{"source":"Home Screen","target":"Cyber Monday Sale Items List","value":97652},{"source":"Order History","target":"Order Status","value":9782},{"source":"Order History","target":"Support","value":2837},{"source":"List Of Items","target":"Sales Offer","value":14678},{"source":"List Of Items","target":"View Item Detail","value":23654},{"source":"Cyber Monday Sale Items List","target":"View Item Detail","value":43889},{"source":"Cyber Monday Sale Items List","target":"Explore Discounts","value":34681},{"source":"Sales Offer","target":"View Item Images","value":12055},{"source":"View Item Detail","target":"View Item Images","value":16793},{"source":"View Item Detail","target":"Add To Cart","value":11537},{"source":"Add To Cart","target":"Pay","value":10144},{"source":"Add To Cart","target":"Explore Discounts","value":4007}]}`

	data2 := `{"nodes":[{"id":"Home Screen","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Order History","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Order Status","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Support","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"List Of Items","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Sales Offer","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"View Item Images","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"View Item Detail","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Cyber Monday Sale Items List","nodeColor":"hsl(0, 72%, 51%)","issues":{"crashes":[{"title":"NullPointerException.java","count":32893},{"title":"LayoutInflaterException.java","count":12874}],"anrs":[{"title":"CyberMondayActivity.java","count":77321},{"title":"CyberMondayFragment.kt","count":6305}]}},{"id":"Add To Cart","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Pay","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}},{"id":"Explore Discounts","nodeColor":"hsl(142, 69%, 58%)","issues":{"crashes":[],"anrs":[]}}],"links":[{"source":"Home Screen","target":"Order History","value":60000},{"source":"Home Screen","target":"List Of Items","value":53356},{"source":"Home Screen","target":"Cyber Monday Sale Items List","value":96652},{"source":"Order History","target":"Order Status","value":9822},{"source":"Order History","target":"Support","value":2287},{"source":"List Of Items","target":"Sales Offer","value":12628},{"source":"List Of Items","target":"View Item Detail","value":53254},{"source":"Cyber Monday Sale Items List","target":"View Item Detail","value":43889},{"source":"Cyber Monday Sale Items List","target":"Explore Discounts","value":34681},{"source":"Sales Offer","target":"View Item Images","value":12055},{"source":"View Item Detail","target":"View Item Images","value":12793},{"source":"View Item Detail","target":"Add To Cart","value":16537},{"source":"Add To Cart","target":"Pay","value":10144},{"source":"Add To Cart","target":"Explore Discounts","value":3007}]}`

	var data string
	randomInt := rand.Intn(100)
	if randomInt > 70 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API server is experiencing intermittent issues"})
		return
	}
	if randomInt%2 == 0 {
		data = data1
	} else {
		data = data2
	}

	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getAppMetrics(c *gin.Context) {
	var af AppFilter

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		msg := `id invalid or missing`
		fmt.Println(msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	af.AppID = id

	if err := c.ShouldBindQuery(&af); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := af.validate(); err != nil {
		msg := "app journey request validation failed"
		fmt.Println(msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg, "details": err.Error()})
		return
	}

	if !af.hasTimeRange() {
		af.setDefaultTimeRange()
	}

	if !af.hasVersion() {
		af.setDefaultVersion()
	}

	// fmt.Println("journey request app id", af.AppID)
	// fmt.Println("journey request from", af.From)
	// fmt.Println("journey request to", af.To)
	// fmt.Println("journey request version", af.Version)

	data1 := `{"adoption":{"users":40000,"totalUsers":200000,"value":20},"app_size":{"value":20,"delta":3.18},"crash_free_users":{"value":98.5,"delta":0.73},"perceived_crash_free_users":{"value":91.3,"delta":-0.51},"multiple_crash_free_users":{"value":76.37,"delta":0.62},"anr_free_users":{"value":98.5,"delta":0.73},"perceived_anr_free_users":{"value":91.3,"delta":0.27},"multiple_anr_free_users":{"value":97.88,"delta":-3.13},"app_cold_launch":{"value":937,"delta":34},"app_warm_launch":{"value":600,"delta":-87},"app_hot_launch":{"value":250,"delta":-55}}`
	data2 := `{"adoption":{"users":49000,"totalUsers":200000,"value":28},"app_size":{"value":20,"delta":3.18},"crash_free_users":{"value":98.2,"delta":0.71},"perceived_crash_free_users":{"value":92.8,"delta":-0.81},"multiple_crash_free_users":{"value":75.49,"delta":0.38},"anr_free_users":{"value":98.3,"delta":0.43},"perceived_anr_free_users":{"value":91.9,"delta":0.77},"multiple_anr_free_users":{"value":97.26,"delta":-2.85},"app_cold_launch":{"value":900,"delta":-200},"app_warm_launch":{"value":600,"delta":-127},"app_hot_launch":{"value":300,"delta":-50}}`

	var data string
	randomInt := rand.Intn(100)
	if randomInt > 70 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API server is experiencing intermittent issues"})
		return
	}
	if randomInt%2 == 0 {
		data = data1
	} else {
		data = data2
	}

	c.Data(http.StatusOK, "application/json", []byte(data))
}

func getAppFilters(c *gin.Context) {
	appId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		msg := `app id invalid or missing`
		fmt.Println(msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	appFiltersMap := map[string]string{
		"59ba1c7f-2a42-4b7f-b9cb-735d25146675": `{
			"version": ["13.2.1", "13.2.2", "13.3.7"],
			"country": [
			  {
				"code": "IN",
				"name": "India"
			  },
			  {
				"code": "CN",
				"name": "China"
			  },
			  {
				"code": "US",
				"name": "USA"
			  }
			],
			"network_provider": ["Airtel", "Jio", "Vodafone"],
			"network_type": ["WiFi", "5G", "4G", "3G", "2G", "Edge"],
			"locale": ["en_IN", "en_US", "en_UK"],
			"device_manufacturer": ["Samsung", "Huawei", "Motorola"],
			"device_name": ["Samsung Galaxy Note 2", "Motorola Razor V2", "Huawei P30 Pro"]
		  }`,
		"243f3214-0f41-4361-8ef3-21d8f5d99a70": `{
			"version": ["13.2.1", "13.2.2", "13.3.9", "13.3.10"],
			"country": [
			  {
				"code": "IN",
				"name": "India"
			  },
			  {
				"code": "CN",
				"name": "China"
			  },
			  {
				"code": "US",
				"name": "USA"
			  }
			],
			"network_provider": ["Airtel", "Jio", "Vodafone"],
			"network_type": ["WiFi", "4G", "3G", "2G"],
			"locale": ["en_IN", "en_US", "en_UK", "zh_HK"],
			"device_manufacturer": ["Samsung", "Huawei", "Motorola"],
			"device_name": ["Samsung Galaxy Note 2", "Motorola Razor V2", "Huawei P30 Pro"]
		  }`,
		"bae4fb9e-07cd-4435-a42e-d99986830c2c": `{
			"version": ["13.3.9", "13.3.12"],
			"country": [
			  {
				"code": "IN",
				"name": "India"
			  },
			  {
				"code": "CN",
				"name": "China"
			  }
			],
			"network_provider": ["Airtel", "Jio", "Vodafone"],
			"network_type": ["WiFi", "4G", "2G"],
			"locale": ["en_IN", "en_UK", "zh_HK"],
			"device_manufacturer": ["Samsung", "Huawei", "Lenovo"],
			"device_name": ["Samsung Galaxy Note 2", "Lenovo Legion Y90", "Huawei P30 Pro"]
		  }`,
		"c6643110-d3e5-4b1c-bfcc-75b46b52ae79": `{
			"version": ["2.2.1", "2.3.2", "3.1.0"],
			"country": [
			  {
				"code": "IT",
				"name": "Italy"
			  },
			  {
				"code": "CN",
				"name": "China"
			  },
			  {
				"code": "DE",
				"name": "Germany"
			  }
			],
			"network_provider": ["Airtel", "Jio", "Vodafone", "Wind Tre", "TIM", "FASTWEB", "Iliad", "Telekom", "O₂"],
			"network_type": ["WiFi", "5G", "4G", "3G", "2G", "Edge"],
			"locale": ["en_IN", "en_US", "en_UK", "it-IT", "en_HK", "zh_HK", "zh_Hans_HK", "ii_CN", "de_DE", "en_DE", "dsb_DE", "hsb_DE"],
			"device_manufacturer": ["Vivo", "Honor", "Oppo", "Huawei", "Xiaomi", "Lenovo", "Samsung"],
			"device_name": ["Samsung Galaxy Note 2", "Motorola Razor V2", "Huawei P30 Pro", "Vivo X900 Pro+", "Oppo Find X6 Pro", "Honor Magic5 Pro", "Xiaomi 13 Pro", "Lenovo Legion Y90", "Samsung Galaxy S23 Ultra"]
		  }`,
		"e2abe28a-f6bc-4f57-88fe-81f10d1c5afc": `{
			"version": ["2.2.1", "2.3.2", "3.1.4"],
			"country": [
			  {
				"code": "IT",
				"name": "Italy"
			  },
			  {
				"code": "CN",
				"name": "China"
			  },
			  {
				"code": "DE",
				"name": "Germany"
			  }
			],
			"network_provider": ["Vodafone", "Wind Tre", "TIM", "FASTWEB", "Iliad", "Telekom", "O₂"],
			"network_type": ["WiFi", "5G", "4G", "3G"],
			"locale": ["en_IN", "en_US", "it-IT", "en_HK", "zh_HK", "zh_Hans_HK", "ii_CN", "de_DE", "en_DE", "dsb_DE", "hsb_DE"],
			"device_manufacturer": ["Vivo", "Honor", "Oppo", "Xiaomi", "Lenovo", "Samsung"],
			"device_name": ["Samsung Galaxy Note 2", "Vivo X900 Pro+", "Oppo Find X6 Pro", "Honor Magic5 Pro", "Xiaomi 13 Pro", "Lenovo Legion Y90", "Samsung Galaxy S23 Ultra"]
		  }`,
		"b17f7003-4ab6-4b1a-a5d8-ed5a72cb4569": `{
			"version": ["2.2.1", "2.3.2", "3.1.1"],
			"country": [
			  {
				"code": "CN",
				"name": "China"
			  },
			  {
				"code": "DE",
				"name": "Germany"
			  }
			],
			"network_provider": ["Vodafone", "Wind Tre", "TIM", "FASTWEB", "Telekom", "O₂"],
			"network_type": ["WiFi", "4G", "3G", "2G", "Edge"],
			"locale": ["en_IN", "en_US", "en_UK", "it-IT", "en_HK", "zh_HK", "zh_Hans_HK", "ii_CN", "de_DE", "en_DE", "dsb_DE", "hsb_DE"],
			"device_manufacturer": ["Vivo", "Honor", "Oppo", "Huawei", "Lenovo", "Samsung"],
			"device_name": ["Samsung Galaxy Note 2", "Motorola Razor V2", "Huawei P30 Pro", "Vivo X900 Pro+", "Oppo Find X6 Pro", "Honor Magic5 Pro", "Lenovo Legion Y90", "Samsung Galaxy S23 Ultra"]
		  }`,
		"20014be8-aaa9-4e56-8810-9f1a48ec1099": `{
			"version": ["3.2.1", "3.3.2", "3.5.1", "3.5.2"],
			"country": [
			  {
				"code": "EE",
				"name": "Estonia"
			  },
			  {
				"code": "CH",
				"name": "Switzerland"
			  }
			],
			"network_provider": ["Telia", "Elisa", "Tele2", "Swisscom", "Sunrise", "Salt"],
			"network_type": ["WiFi", "5G", "4G", "3G", "2G", "Edge"],
			"locale": ["en_IN", "en_US", "en_CA", "it-CH", "rm_CH", "gsw_CH", "et_EE"],
			"device_manufacturer": ["Vivo", "Honor", "Oppo", "Huawei", "Lenovo", "Samsung", "Realme", "OnePlus", "Nokia"],
			"device_name": ["Huawei P30 Pro", "OnePlus 11 Pro", "Vivo X900 Pro+", "Noia X30", "Oppo Find X6 Pro", "Honor Magic5 Pro", "Lenovo Legion Y90", "Samsung Galaxy S23 Ultra"]
		  }`,
		"463c959c-94c2-4f49-bd2b-6caab360c152": `{
			"version": ["3.3.2", "3.5.1", "3.5.2", "3.5.3"],
			"country": [
			  {
				"code": "EE",
				"name": "Estonia"
			  },
			  {
				"code": "CH",
				"name": "Switzerland"
			  }
			],
			"network_provider": ["Telia", "Elisa", "Tele2", "Swisscom", "Sunrise", "Salt"],
			"network_type": ["5G", "4G", "3G", "2G"],
			"locale": ["en_CA", "it-CH", "rm_CH", "gsw_CH", "et_EE"],
			"device_manufacturer": ["Honor", "Oppo", "Huawei", "Lenovo", "Samsung", "Realme", "OnePlus", "Nokia"],
			"device_name": ["Huawei P30 Pro", "OnePlus 11 Pro", "Noia X30", "Oppo Find X6 Pro", "Honor Magic5 Pro", "Lenovo Legion Y90", "Samsung Galaxy S23 Ultra"]
		  }`,
		"2a7f230e-6d5e-4036-b4e6-1102c22f4433": `{
			"version": ["3.5.1", "3.5.2", "3.5.3", "3.6.0"],
			"country": [
			  {
				"code": "EE",
				"name": "Estonia"
			  },
			  {
				"code": "CH",
				"name": "Switzerland"
			  }
			],
			"network_provider": ["Telia", "Elisa", "Tele2", "Swisscom", "Sunrise", "Salt"],
			"network_type": ["WiFi", "5G", "4G", "2G", "Edge"],
			"locale": ["en_IN", "en_US", "en_CA", "it-CH", "rm_CH", "gsw_CH", "et_EE"],
			"device_manufacturer": ["Vivo", "Honor", "Oppo", "Huawei", "Realme", "OnePlus", "Nokia"],
			"device_name": ["Huawei P30 Pro", "OnePlus 11 Pro", "Vivo X900 Pro+", "Noia X30", "Oppo Find X6 Pro", "Honor Magic5 Pro"]
		  }`,
	}

	appFilters := appFiltersMap[appId.String()]

	if appFilters == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no app filters exists for app [%s]", appId.String())})
	} else {
		c.Data(http.StatusOK, "application/json", []byte(appFilters))
	}

}