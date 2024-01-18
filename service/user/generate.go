package main

import (
	"math/rand"
)

var (
	avatars = map[int]string{
		0: "20220612164733_72d8b.jpg",
		1: "20220622180647_d4cb5.jpg",
		2: "20220709150824_97667.jpg",
		3: "20220801091937_fc599.jpg",
		4: "20220801091938_56fad.jpg",
		5: "20220801091938_9f3d5.jpg",
		6: "20220801091939_9a475.jpg",
		7: "20220801165632_07749.jpg",
		8: "20220801204306_f3f98.jpg",
		9: "20220906115559_aff77.jpg",
	}
	backgrounds = map[int]string{
		0:  "background0.jpg",
		1:  "background1.jpg",
		2:  "background2.jpg",
		3:  "background3.jpg",
		4:  "background4.jpg",
		5:  "background5.jpg",
		6:  "background6.jpg",
		7:  "background7.jpg",
		8:  "background8.jpg",
		9:  "background9.jpg",
		10: "background10.jpg",
		11: "background11.jpg",
		12: "background12.jpg",
	}
	signatrues = map[int]string{
		0:  "夜猫子协会常任理事",
		1:  "赖床锦标赛冠军得主",
		2:  "深夜搞颜色积极分子",
		3:  "贫困大赛形象代言人",
		4:  "魔仙堡废话冠军",
		5:  "迪士尼在逃保洁阿姨。",
		6:  "非官方认证平平无奇说废话小天才",
		7:  "中央戏精学院教授",
		8:  "口吐芬芳专业教授",
		9:  "顶级外卖鉴赏师",
		10: "秃头选拔赛形象大使",
		11: "互联网冲浪金牌选手",
		12: "国家一级退堂鼓选手",
		13: "国家一级抬杠运动员",
		14: "耳机依赖患者",
		15: "宇宙一级潜在鸽王",
		16: "退役熬夜选手",
		17: "拖延俱乐部顶级VIP",
		18: "退役魔法少女",
		19: "脆皮鸭文学爱好者",
		20: "2023年广东省高考状元老乡",
		21: "铠甲勇士赵本山",
	}
)

func generateAvatar() string {
	return avatars[rand.Intn(len(avatars))]
}

func generateImage() string {
	return backgrounds[rand.Intn(len(backgrounds))]
}

func generateSignatrue() string {
	return signatrues[rand.Intn(len(signatrues))]
}
