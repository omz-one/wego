package payment

import (
	"github.com/godcong/wego/core"
	"github.com/godcong/wego/util"
)

/*Merchant 账单 */
type Merchant struct {
	*Payment
}

func newMerchant(p *Payment) *Merchant {
	return &Merchant{
		Payment: p,
	}
}

/*NewMerchant 账单 */
func NewMerchant(config *core.Config) *Merchant {
	return newMerchant(NewPayment(config))
}

func (m *Merchant) AddSubMerchant(maps util.Map) core.Response {
	return m.manage("add", maps)
}

func (m *Merchant) QuerySubMerchantByMerchantId(id string) core.Response {
	return m.manage("query", util.Map{"micro_mch_id": id})
}

func (m *Merchant) QuerySubMerchantByWeChatId(id string) core.Response {
	return m.manage("query", util.Map{"recipient_wechatid": id})
}

func (m *Merchant) ModifyInfo(maps util.Map) {

}

func (m *Merchant) manage(action string, maps util.Map) core.Response {

	maps.Join(util.Map{
		"appid":      m.GetString("app_id"),
		"nonce_str":  "",
		"sub_mch_id": "",
		"sub_appid":  "",
	})
	params := util.Map{
		core.DataTypeXML:      m.initRequest(maps),
		core.DataTypeQuery:    util.Map{"action": action},
		core.DataTypeSecurity: m.Config,
	}
	return m.client.Request(Link(mchSubmchmanage), "post", params)
}
