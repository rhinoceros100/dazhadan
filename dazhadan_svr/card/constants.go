package card

const (
	CardType_Diamond	int = iota + 1 		//方片
	CardType_Club	 			        //梅花
	CardType_Heart		 			//红桃
	CardType_Spade					//黑桃
	CardType_BlackJoker				//小王
	CardType_RedJoker				//大王
)

const (
	CardsType_NO		int = iota		//没有任何牌型
	CardsType_SINGLE				//单牌
	CardsType_STAIGHT				//顺子
	CardsType_PAIR			           	//对子
	CardsType_PAIRS		           		//连对
	CardsType_32		           		//三带二 5
	CardsType_43		           		//四带三
	CardsType_54		           		//五带四
	CardsType_65		           		//六带五
	CardsType_76		           		//七带六
	CardsType_87		           		//八带七
	CardsType_PLANE32	           		//三带二飞机 11
	CardsType_PLANE43	           		//四带三飞机
	CardsType_PLANE54	           		//五带四飞机
	CardsType_PLANE65	           		//六带五飞机
	CardsType_PLANE76	           		//七带六飞机
	CardsType_PLANE87	           		//八带七飞机
	CardsType_510K		= 21      		//五十K
	CardsType_TRUE510K	= 22          		//真五十K
	CardsType_BOMB4		= 23           		//四炸
	CardsType_BOMB5		= 24	           	//五炸
	CardsType_BOMB6		= 25        		//六炸
	CardsType_BOMB7		= 26       		//七炸
	CardsType_BOMB8		= 27       		//八炸
	CardsType_JOKER		= 28       		//王炸
)

const (
	InfoType_Normal		int32 = iota		//普通结果
	InfoType_Shuangji				//双基
	InfoType_PlayAloneSucc				//打独成功
	InfoType_PlayAloneFail				//打独失败
)
