package object

import (
	"errors"
	"fmt"
)

type ObjectType string

const (
	COLLECTABLE_OBJ  = "COLLECTABLE"
	EXCHANGEABLE_OBJ = "EXCHANGEABLE"
	OFFER_OBJ        = "OFFER"
	TABLEVIEW_OBJ    = "TABLEVIEW"
	ERROR_OBJ        = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type CollectableName string

type Collectable struct {
	Name   CollectableName
	Amount int
	Owner  string
}

func (c *Collectable) Type() ObjectType { return COLLECTABLE_OBJ }
func (c *Collectable) Inspect() string {
	return fmt.Sprintf("(%s, Name: %q, Amount: %d, Owner: %s)", c.Type(), c.Name, c.Amount, c.Owner)
}

type Exchangeable struct {
	Amount int
	Owner  string
}

func (ex *Exchangeable) Type() ObjectType { return EXCHANGEABLE_OBJ }
func (ex *Exchangeable) Inspect() string {
	return fmt.Sprintf("( %s, Amount: %d, owner: %s)", ex.Type(), ex.Amount, ex.Owner)
}

type Offer struct {
	Id            int
	LCollectables []*Collectable
	RCollectables []*Collectable
}

func (of *Offer) Type() ObjectType { return OFFER_OBJ }
func (of *Offer) Inspect() string {
	lcollectables := []string{}
	for _, collectable := range of.LCollectables {
		lcollectables = append(lcollectables, collectable.Inspect())
	}
	rcollectables := []string{}
	for _, collectable := range of.RCollectables {
		rcollectables = append(rcollectables, collectable.Inspect())
	}
	return fmt.Sprintf("(%s, Id: %d,\n\tLCollectable: [%s\t],\n\tRCollectable:[%s\t])", of.Type(), of.Id, lcollectables, rcollectables)
}

type ViewObject struct {
	Collections []*Collectable
}

func (vo *ViewObject) Type() ObjectType { return TABLEVIEW_OBJ }
func (vo *ViewObject) Inspect() string {
	msg := fmt.Sprintf("| %-15s | %-8s | %-15s |\n", "Name", "Amount", "Owner")

	for _, tv := range vo.Collections {
		// Aplicamos los mismos anchos a los datos de la fila
		msg += fmt.Sprintf("| %-15s | %-8d | %-15s |\n", tv.Name, tv.Amount, tv.Owner)
	}
	return msg
}

type Identifier struct {
	Value CollectableName
}

func (i *Identifier) Type() ObjectType { return TABLEVIEW_OBJ }
func (i *Identifier) Inspect() string  { return fmt.Sprintf("(%s, Value: %s)", i.Type(), i.Value) }

const (
	// Argentina (AR)
	AR_LM10 = "AR-LM10" // Lionel Messi
	AR_AD11 = "AR-AD11" // Ángel Di María
	AR_EM23 = "AR-EM23" // Emiliano Martínez
	AR_RD7  = "AR-RD7"  // Rodrigo De Paul
	AR_JA9  = "AR-JA9"  // Julián Álvarez

	// Brasil (BR)
	BR_VJ7  = "BR-VJ7"  // Vinícius Júnior
	BR_NJ10 = "BR-NJ10" // Neymar Jr.
	BR_AB1  = "BR-AB1"  // Alisson Becker
	BR_RG11 = "BR-RG11" // Rodrygo Goes
	BR_MQ5  = "BR-MQ5"  // Marquinhos

	// Francia (FR)
	FR_KM10 = "FR-KM10" // Kylian Mbappé
	FR_AG7  = "FR-AG7"  // Antoine Griezmann
	FR_OG9  = "FR-OG9"  // Olivier Giroud
	FR_NK13 = "FR-NK13" // N'Golo Kanté
	FR_AT8  = "FR-AT8"  // Aurélien Tchouaméni

	// España (ES)
	ES_AM7  = "ES-AM7"  // Álvaro Morata
	ES_RH16 = "ES-RH16" // Rodrigo Hernández (Rodri)
	ES_DC2  = "ES-DC2"  // Dani Carvajal
	ES_LY19 = "ES-LY19" // Lamine Yamal
	ES_PG20 = "ES-PG20" // Pedri González

	// Reino Unido / Inglaterra (GB)
	GB_HK9  = "GB-HK9"  // Harry Kane
	GB_JB10 = "GB-JB10" // Jude Bellingham
	GB_BS7  = "GB-BS7"  // Bukayo Saka
	GB_PF11 = "GB-PF11" // Phil Foden
	GB_DR4  = "GB-DR4"  // Declan Rice
)

var validPlayerCodesSet = map[CollectableName]struct{}{
	// Argentina (AR)
	AR_LM10: {},
	AR_AD11: {},
	AR_EM23: {},
	AR_RD7:  {},
	AR_JA9:  {},

	// Brasil (BR)
	BR_VJ7:  {},
	BR_NJ10: {},
	BR_AB1:  {},
	BR_RG11: {},
	BR_MQ5:  {},

	// Francia (FR)
	FR_KM10: {},
	FR_AG7:  {},
	FR_OG9:  {},
	FR_NK13: {},
	FR_AT8:  {},

	// España (ES)
	ES_AM7:  {},
	ES_RH16: {},
	ES_DC2:  {},
	ES_LY19: {},
	ES_PG20: {},

	// Reino Unido / Inglaterra (GB)
	GB_HK9:  {},
	GB_JB10: {},
	GB_BS7:  {},
	GB_PF11: {},
	GB_DR4:  {},
}

func (p CollectableName) IsValid() bool {
	_, exists := validPlayerCodesSet[p]
	return exists
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Creamos una base de datos en memoria. Mejora: JSON
type Environment struct {
	users []struct {
		username string
		password string
	} //[{"username":"hola","password":"hola"}]
	collectables map[string][]*Collectable //{"username":[]{"collectableName","Amount"}}
	exchangeable map[string][]*Collectable //{"username":[]{"collectableName","amount"}}
	offers       map[string][]*Offer       //{"username":[]{"id",lcollectables,rcollectables}}

	//Usuario actual
	actualUser string
}

func NewEnvironment() *Environment {
	//Creacion de la DB.
	users := []struct {
		username string
		password string
	}{
		{"pepe", "123456"},
		{"pancho", "pancho"},
	}

	collectables := make(map[string][]*Collectable)

	collectables["pepe"] = append(collectables["pepe"],
		&Collectable{Name: AR_LM10, Amount: 5, Owner: "pepe"},
		&Collectable{Name: BR_VJ7, Amount: 3, Owner: "pepe"},
		&Collectable{Name: FR_KM10, Amount: 12, Owner: "pepe"},
		&Collectable{Name: ES_LY19, Amount: 2, Owner: "pepe"},
		&Collectable{Name: GB_JB10, Amount: 7, Owner: "pepe"},
	)

	collectables["pancho"] = append(collectables["pancho"],
		&Collectable{Name: AR_EM23, Amount: 10, Owner: "pancho"},
		&Collectable{Name: BR_NJ10, Amount: 4, Owner: "pancho"},
		&Collectable{Name: FR_AG7, Amount: 8, Owner: "pancho"},
		&Collectable{Name: ES_RH16, Amount: 6, Owner: "pancho"},
		&Collectable{Name: GB_HK9, Amount: 15, Owner: "pancho"},
	)
	exchangeable := make(map[string][]*Collectable)
	offers := make(map[string][]*Offer)

	//Creación de una oferta de prueba.
	// offers["pepe"] = append(offers["pepe"], &Offer{
	// 	Id: 1,
	// 	LCollectables: []*Collectable{
	// 		{Name: AR_LM10, Amount: 2},
	// 		{Name: AR_AD11, Amount: 1},
	// 	},
	// 	RCollectables: []*Collectable{
	// 		{Name: FR_KM10, Amount: 1},
	// 		{Name: BR_VJ7, Amount: 3},
	// 		{Name: ES_LY19, Amount: 2},
	// 	},
	// })

	//Creación de exchangeables de prueba
	exchangeable["pepe"] = append(exchangeable["pepe"],
		&Collectable{Name: AR_LM10, Amount: 5, Owner: "pepe"},
		&Collectable{Name: BR_VJ7, Amount: 3, Owner: "pepe"},
		&Collectable{Name: FR_KM10, Amount: 12, Owner: "pepe"},
		&Collectable{Name: ES_LY19, Amount: 2, Owner: "pepe"},
		&Collectable{Name: GB_JB10, Amount: 7, Owner: "pepe"},
	)

	exchangeable["pancho"] = append(exchangeable["pancho"],
		&Collectable{Name: AR_EM23, Amount: 10, Owner: "pancho"},
		&Collectable{Name: BR_NJ10, Amount: 4, Owner: "pancho"},
		&Collectable{Name: FR_AG7, Amount: 8, Owner: "pancho"},
		&Collectable{Name: ES_RH16, Amount: 6, Owner: "pancho"},
		&Collectable{Name: GB_HK9, Amount: 15, Owner: "pancho"},
	)
	return &Environment{users, collectables, exchangeable, offers, "pepe"}
}

func (env *Environment) GetCollectables() []*Collectable {
	return env.collectables[env.actualUser]
}

func (env *Environment) SetExchangeableCollection(queryCollectables []*Collectable) error {
	dbCollectables := env.GetCollectables()
	userExchangeableCollections, ok := env.exchangeable[env.actualUser]

	if len(dbCollectables) == 0 {
		return errors.New("no stock: No tienes más stock de coleccionables.")
	}

	if !ok {
		userExchangeableCollections = []*Collectable{}
	}

	for _, dc := range dbCollectables {
		for _, qc := range queryCollectables {
			if !qc.Name.IsValid() {
				return fmt.Errorf("unknown collectable: El coleccionable %s no existe.", qc.Name)
			}
			if dc.Name == qc.Name {
				if qc.Amount > dc.Amount {
					return fmt.Errorf("no stock: No tienes suficiente coleccionables %s para ofrecer. (Tienes %d)", qc.Name, dc.Amount) //La query me pide coleccionables que no tengo
				}
				dc.Amount -= qc.Amount
				found := false

				//Nos fijamos si ya existe
				for _, ex := range userExchangeableCollections {
					if ex.Name != qc.Name {
						continue
					}
					found = true
					ex.Amount += qc.Amount
				}
				//Si no existe agregarlos al intercambiable.
				if !found {
					env.exchangeable[env.actualUser] = append(env.exchangeable[env.actualUser], qc)
				}
			}
		}
	}

	//Actualizamos la db
	env.collectables[env.actualUser] = dbCollectables
	env.exchangeable[env.actualUser] = userExchangeableCollections
	return nil
}

func (env *Environment) GetActualUser() string {
	return env.actualUser
}

func (env *Environment) GetExchangeables(collectionName CollectableName) *ViewObject {
	var c []*Collectable
	for user, collection := range env.exchangeable {
		if user == env.actualUser {
			continue
		}
		for _, item := range collection {
			if item.Name == collectionName {
				c = append(c, item)
			}
		}
	}
	return &ViewObject{Collections: c}
}
