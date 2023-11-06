package util

import "math/rand"

var (
	avatars = map[int]string{
		0: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220612164733_72d8b.jpg",
		1: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220622180647_d4cb5.jpg",
		2: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220709150824_97667.jpg",
		3: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220801091937_fc599.jpg",
		4: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220801091938_56fad.jpg",
		5: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220801091938_9f3d5.jpg",
		6: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220801091939_9a475.jpg",
		7: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220801165632_07749.jpg",
		8: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220801204306_f3f98.jpg",
		9: "http://s2a5yl4lg.hn-bkt.clouddn.com/20220906115559_aff77.jpg",
	}
	backgrounds = map[int]string{
		0:  "https://img2.wallspic.com/previews/0/6/1/6/7/176160/176160-eren_yeager-gong_ji_de_ju_ren-a_mingarlert-yi_shu-hai_bao-500x.jpg",
		1:  "https://img1.wallspic.com/previews/1/4/1/6/7/176141/176141-ban_yuan_ding-bing_chuan_dian-yue_sai_mi_di_shan_gu-zi_ran_jing_guan-ren_men_zai_zi_ran_jie-500x.jpg",
		2:  "https://img3.wallspic.com/previews/9/1/0/6/7/176019/176019-wu_daotouka-kenkaneki-dong_jing_shi_shi_gui-hai_bao-azure-500x.jpg",
		3:  "https://img3.wallspic.com/previews/3/8/9/5/7/175983/175983-chu_yin_wei_lai-zui-wei_xiao-ka_tong-azure-500x.jpg",
		4:  "https://img1.wallspic.com/previews/8/5/8/5/7/175858/175858-qi_fen-yu_hui-zi_ran_jing_guan-yang_guang-ji_yun-500x.jpg",
		5:  "https://img3.wallspic.com/previews/9/9/7/5/7/175799/175799-playstation2-ps4-ps3you_xi_ji-playstation_5-playstation_shang_dian-500x.jpg",
		6:  "https://img3.wallspic.com/previews/4/9/7/4/7/174794/174794-di_si_ni_dian_ying-di_shi_ni_gong_si-xiu_xian-le_qu-yu_le-500x.jpg",
		7:  "https://img3.wallspic.com/previews/0/7/0/5/7/175070/175070-zui-ka_tong-hen_ku_de-yi_shu-dian_lan_se_de-500x.jpg",
		8:  "https://img1.wallspic.com/previews/6/4/2/5/7/175246/175246-ping_guo_bkc_meng_mai-ping_guo-air-huang_se_de-dian_lan_se_de-500x.jpg",
		9:  "https://img1.wallspic.com/previews/5/8/0/5/7/175085/175085-yi_gexbox-huo_ying_ren_zhe-ka_tong-yi_shu-dong_hua-500x.jpg",
		10: "https://img3.wallspic.com/previews/4/9/9/4/7/174994/174994-yuan_quan-jian_tie_hua_de-yi_shu-pin_hong_se-er_tong_yi_shu-500x.jpg",
		11: "https://img2.wallspic.com/previews/5/4/6/4/7/174645/174645-hong_hu_li-fu_ke_si-bei_ji_hu-xiao_lu-lu_de_dong_wu-500x.jpg",
		12: "https://img3.wallspic.com/previews/2/3/6/4/7/174632/174632-yi_shu-yi_shu_zhan-xian_dai_yi_shu-you_hua-azure-500x.jpg",
	}
	signatrues = map[int]string{
		0: "夜猫子协会常任理事",
		1:    "赖床锦标赛冠军得主",
		2:    "深夜搞颜色积极分子",
		3:    "贫困大赛形象代言人",
		4:    "魔仙堡废话冠军",
		5:    "迪士尼在逃保洁阿姨。",
		6:    "非官方认证平平无奇说废话小天才",
		7:    "中央戏精学院教授",
		8:    "口吐芬芳专业教授",
		9:    "顶级外卖鉴赏师",
		10:   "秃头选拔赛形象大使",
		11:   "互联网冲浪金牌选手",
		12:   "国家一级退堂鼓选手",
		13:   "国家一级抬杠运动员",
		14:   "耳机依赖患者",
		15:   "宇宙一级潜在鸽王",
		16:   "退役熬夜选手",
		17:   "拖延俱乐部顶级VIP",
		18:   "退役魔法少女",
		19:   "脆皮鸭文学爱好者",
		20:   "2023年广东省高考状元老乡",
		21:   "铠甲勇士赵本山",
	}
)

func GenerateAvatar() string {
	return avatars[rand.Intn(len(avatars))]
}

func GenerateImage() string {
	return backgrounds[rand.Intn(len(backgrounds))]
}

func GenerateSignatrue() string {
	return signatrues[rand.Intn(len(signatrues))]
}
