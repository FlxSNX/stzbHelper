<script setup>
import { NButton,NH1,NDivider,darkTheme,NConfigProvider,NGlobalStyle,NCard,NForm,NFormItem,NInput,NSpin,NSelect,NSwitch,NFormItemGi } from 'naive-ui';
import { reportList } from '../api';
import { onMounted,onBeforeUnmount,ref } from 'vue';
import { herocfg as hcfg, skillcfg as scfg } from '../cfg';
import { useRouter } from 'vue-router'

const router = useRouter()

const herocfg = JSON.parse(hcfg);
const skillcfg = JSON.parse(scfg);

console.log(herocfg[100762].iconId);

const report = ref([]);
const nextid = ref(0);
const atkname = ref('');
const atkunionname = ref('');
const atkhp = ref('');
const atkstar = ref('');
const atklevel = ref('');
const loaddataing = ref(false);
const reportTotal = ref(0);
const serachType = ref(null);
const noviewnpc = ref(false);
const havemore = ref(true);
const serachTypeList = [
	{
		value: '1',
		label: '双方其一满足条件'
	},
	{
		value: '2',
		label: '只筛选进攻方'
	},
	{
		value: '3',
		label: '只筛选防守方'
	},
	{
		value: '4',
		label: '双方都需要满足条件(名字和同盟名字除外)'
	},
];

const resetReport = () => {
	atkname.value = '';
	atkstar.value = '';
	atklevel.value = '';
	serachType.value = null;
	noviewnpc.value = false;
	atkunionname.value = '';
	atkhp.value = '';
	getReportListData(1);
}

const getrightheroskillindex = (index) => {
	if(index == 0){
		return 5;
	}else if(index == 1){
		return 4;
	}else{
		return 3
	}
}

function formatTimestamp(timestamp) {
    // 将秒级时间戳转换为毫秒级
    const date = new Date(timestamp * 1000);

    // 获取年、月、日、时、分、秒
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从0开始，所以要加1
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    // 拼接成目标格式
    return `${year}/${month}/${day} ${hours}:${minutes}:${seconds}`;
}

function splitwid(num) {
    // 将数字转换为字符串
    const numStr = num.toString();
    
    // 获取后四位
    const lastFour = numStr.slice(-4);
    
    // 获取前面的部分
    const firstPart = numStr.slice(0, -4);
    
    // 将后四位转换为数字，并去掉前导零
    const lastFourNumber = parseInt(lastFour, 10);
    
    // 返回结果
    return `${firstPart},${lastFourNumber}`
}

const getReportListData = (clear=false) => {
	if(clear){
		nextid.value = 0;
		report.value = [];
		havemore.value = true;
		// return
	}
	if(loaddataing.value == true)return;
	loaddataing.value = true;
	reportList({
		nextid:nextid.value,
		atkname:atkname.value,
		atkunionname:atkunionname.value,
		atkhp:atkhp.value,
		atkstar:atkstar.value,
		atklevel:atklevel.value,
		nonpc:noviewnpc.value ? 1 : 0,
		type:serachType.value == null ? 1 : serachType.value
	}).then((v) => {
		// console.log(v);
		if(v.status == 200){
			let resp = v.data
			let data = resp.data.report;
			reportTotal.value = resp.data.total;
			console.log(data);
			if(resp.data.report.length > 0){
				nextid.value = resp.data.report[resp.data.report.length-1].id
				console.log(nextid.value);
				data.forEach(e => {
					console.log(e);
					// e.data = JSON.parse(e.data)
					let all_skill_info_arr = e.all_skill_info.split(';').filter(item => item !== '');
					let all_skill_info = [];
					all_skill_info_arr.forEach(e => {
						e = e.split(',')
						let obj = {
							index:e[0],
							skill1id:e[1],
							skill1level:e[2],
							skill2id:e[3],
							skill2level:e[4],
							skill3id:e[5],
							skill3level:e[6],
						}
						all_skill_info.push(obj)
					});
					let attack_all_hero_info = [];
					let attack_all_hero_info_arr = e.attack_all_hero_info.split(';').filter(item => item !== '');
					// 0,0,0,0,0,0;5,1,1,0,1,0;5,1,1,0,1,19;5,1,1,0,1,0; 满红的数据
					let attack_advance = e.attack_advance.split(';').filter(item => item !== '');
					let attacker_gear_info = e.attacker_gear_info.split(';').filter(item => item !== '');
					let attack_curhp = 0;
					let attack_hurthp = 0;
					attack_all_hero_info_arr.forEach((atk,index) => {
						let atk_arr = atk.split(",")
						if(atk_arr[0].slice(0, 2) === "13")atk_arr[0] =  "10"+atk_arr[0].slice(2)
						attack_all_hero_info.push({
							id:atk_arr[0],
							level:atk_arr[1],
							army:atk_arr[2],
							curArmy:atk_arr[3],
							hurtArmy:atk_arr[4],
							advance:attack_advance[index+1].split(','),
							attacker_gear_info:attacker_gear_info[index+1].split(','),
						});
						attack_curhp += parseInt(atk_arr[3]);
						attack_hurthp += parseInt(atk_arr[3]) + parseInt(atk_arr[4]);
					});
					let defend_all_hero_info = [];
					let defend_all_hero_info_arr = e.defend_all_hero_info.split(';').filter(item => item !== '');
					let defend_advance = e.defend_advance.split(';').filter(item => item !== '');
					let defend_curhp = 0;
					let defend_hurthp = 0;
					defend_all_hero_info_arr.forEach((atk,index) => {
						let atk_arr = atk.split(",")
						if(atk_arr[0].slice(0, 2) === "13")atk_arr[0] =  "10"+atk_arr[0].slice(2)
						let advindex = 0;
						if(index == 0){
							advindex = 2;
						}else if(index == 1){
							advindex = 1
						}
						defend_all_hero_info.push({
							id:atk_arr[0],
							level:atk_arr[1],
							army:atk_arr[2],
							curArmy:atk_arr[3],
							hurtArmy:atk_arr[4],
							advance:defend_advance[advindex].split(',')
						})
						defend_curhp += parseInt(atk_arr[3]);
						defend_hurthp += parseInt(atk_arr[3]) + parseInt(atk_arr[4]);
					});
					e.attack_all_hero_info = attack_all_hero_info;
					e.defend_all_hero_info = defend_all_hero_info;
					e.defend_curhp = defend_curhp;
					e.defend_hurthp = defend_hurthp;
					e.attack_curhp = attack_curhp;
					e.attack_hurthp = attack_hurthp;
					e.attack_hp = e.attack_hp;
					e.defend_hp = e.defend_hp;
					e.all_skill_info = all_skill_info;
					report.value.push(e)
				});
				
				console.log(report.value);
			}else{
				havemore.value = false;
			}
			loaddataing.value = false;
		}
	}).catch(e => {
		loaddataing.value = false;
	});
}

const handleScroll = () => {
	const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
	const scrollHeight = document.documentElement.scrollHeight || document.body.scrollHeight;
	const clientHeight = document.documentElement.clientHeight || document.body.clientHeight;

	if (scrollTop + clientHeight >= scrollHeight - 30) { // 10 是一个缓冲值，可以根据需要调整
		// console.log('已经滚动到页面底部');
		// 在这里可以触发加载更多数据的操作
		if(havemore.value == true)getReportListData();
	}
};

onMounted( _ => {
	getReportListData();
	window.addEventListener('scroll', handleScroll);
})

onBeforeUnmount(() => {
	window.removeEventListener('scroll', handleScroll);
});
</script>

<template>
	<n-config-provider :theme="darkTheme">
		<div style="text-align: center;">
			<!-- <n-h1>率土之滨战报数据站</n-h1>
			<n-divider dashed>
				知彼知己，百战不殆
			</n-divider>
			<n-card title="说明" style="text-align: left;">

				<p>战报数据目前是20秒更新1次, 后续如果有问题可能会调整为3分钟更新1次</p>

				<p>战斗地点显示未知并且坐标X与Y都较小时 有可能是演武和其他不在大地图上的战报</p>

				<p><br>X4472已收集战报数量: {{reportTotal}}</p>
			</n-card> -->
			<!-- <n-radio-group v-model:value="serachType" default-value="1">
				<n-radio-button
					v-for="item in serachTypeList"
					:key="item.value"
					:value="item.value"
					:label="item.label"
					:default-checked="item.value == serachType"
				/>
			</n-radio-group> -->

			<div style="margin-top:16px;width: 30%;">
				<n-select v-model:value="serachType" :options="serachTypeList" placeholder="选择筛选模式 默认双方其一满足条件" />
			</div>
			<n-form
				inline
				:label-width="80"
				style="margin-top: 16px"
			>
				<n-form-item label="名字">
					<n-input placeholder="" v-model:value="atkname"/>
				</n-form-item>
				<n-form-item label="同盟名字">
						<n-input placeholder="" v-model:value="atkunionname" />
				</n-form-item>
				<n-form-item label="红度">
						<n-input placeholder="" v-model:value="atkstar" />
				</n-form-item>
				<n-form-item label="兵力">
						<n-input placeholder="" v-model:value="atkhp" />
				</n-form-item>
				<n-form-item label="等级">
						<n-input placeholder="" v-model:value="atklevel" />
				</n-form-item>
				<n-form-item>
					<n-button attr-type="button" @click="getReportListData(1)">
						筛选
					</n-button>
				</n-form-item>
				<n-form-item>
					<n-button attr-type="button" @click="resetReport(1)">
						重置
					</n-button>
				</n-form-item>
			</n-form>

			<n-form 
				inline
				:label-width="80"
				style="margin-top: 16px"
			>
				<n-form-item :span="12" label="不查看与NPC的战报">
					<n-switch v-model:value="noviewnpc" />
				</n-form-item>
                <n-form-item>
					<n-button attr-type="button" @click="()=>{
                        
                        router.push('/team')
                    }">
						前往查询队伍
					</n-button>
				</n-form-item>
			</n-form>
			<div class="report">
				<div class="report-item" v-for="item in report">
					<div style="align-self: flex-start;margin-left: 8px;">战报ID：{{ item.battle_id }}</div>
					<div class="report-header">
						<div class="left">
							<div class="info">
								<div>{{item.attack_union_name  ? "【盟】"+item.attack_union_name : ""}}</div>
								<div>{{ item.attack_name == "" && item.npc ? "守军" : item.attack_name }}</div>
							</div>
							<div class="report-bar" style="justify-content: flex-end;">
								<div style="margin-right: 8px;">{{item.attack_curhp}}/<span style="color:#8d8d8d">{{ item.attack_hp }}</span></div>
								<div class="bar-box">
									<div class="bar1 left" :style="`width:${item.attack_hurthp / item.attack_hp * 100}%`"></div>
									<div class="bar2 left" :style="`width:${item.attack_curhp / item.attack_hp * 100}%`"></div>
								</div>
							</div>
						</div>
						<div class="center">
							<div>{{ item.wid_name == "" || item.wid_name == " " ? "未知" : item.wid_name }} ({{ splitwid(item.wid) }})</div>
						</div>
						<div class="right">
							<div class="info">
								<div>{{item.defend_union_name ? item.defend_union_name+"【盟】" : ""}}</div>
								<div>{{ item.defend_name == "" && item.npc ? "守军" : item.defend_name }}</div>
							</div>
							<div class="report-bar">
								<div class="bar-box">
									<div class="bar1" :style="`width:${item.defend_hurthp / item.defend_hp * 100}%`"></div>
									<div class="bar2" :style="`width:${item.defend_curhp / item.defend_hp * 100}%`"></div>
								</div>
								<div style="margin-left: 8px;">{{ item.defend_curhp }}/<span style="color:#8d8d8d">{{ item.defend_hp }}</span></div>
							</div>
						</div>
					</div>
					<div class="report-content">
						<div class="report-hero-left">
							<div v-for="(hero,heroindex) in item.attack_all_hero_info">
								<div v-if="hero.id != 0">{{ herocfg[hero.id]?.name ? herocfg[hero.id]?.name : hero.id }}</div>
								<div class="report-hero" v-if="hero.id != 0">
									
									<div class="hero" style="height: 140px;">
										<img :src="`https://g0.gph.netease.com/ngsocial/community/stzb/cn/cards/cut/card_medium_${herocfg[hero.id]?.iconId ? herocfg[hero.id]?.iconId : hero.id }.jpg?gameid=g10`">
										<div class="bottom">
											<div>LV.{{ hero.level }}</div>
											<div>{{ hero.army }}</div>
										</div>
										<div class="stars">
											<!-- <div class="star" :class="{'up': item.Up >= (i+1)}" v-for="(s,i) in HEROS[item.Id].star"></div> -->
											<div class="star" :class="{'up': hero.advance[0] >= (i+1)}" v-for="(s,i) in herocfg[hero.id]?.quality ? herocfg[hero.id]?.quality : 0"></div>
										</div>
									</div>
									<!-- <div>{{ hero.attacker_gear_info[0] }} {{ hero.attacker_gear_info[1] }} {{ hero.attacker_gear_info[2] }}</div> -->
								</div>	
								<div class="skills" v-if="hero.id != 0 && item.all_skill_info.length > 0">
									<div class="skill" v-if="item.all_skill_info[heroindex]?.skill1id != 0">
										{{ skillcfg[item.all_skill_info[heroindex]?.skill1id]?.name ? skillcfg[item.all_skill_info[heroindex]?.skill1id]?.name : '未知' }} 
										[{{ item.all_skill_info[heroindex]?.skill1level}}级]
									</div>
									<div class="skill" v-if="item.all_skill_info[heroindex]?.skill2id != 0">
										{{ skillcfg[item.all_skill_info[heroindex]?.skill2id]?.name ? skillcfg[item.all_skill_info[heroindex]?.skill2id]?.name : '未知' }} 
										[{{ item.all_skill_info[heroindex]?.skill2level}}级]
									</div>
									<div class="skill" v-if="item.all_skill_info[heroindex]?.skill3id != 0">
										{{ skillcfg[item.all_skill_info[heroindex]?.skill3id]?.name ? skillcfg[item.all_skill_info[heroindex]?.skill3id]?.name : '未知' }} 
										[{{ item.all_skill_info[heroindex]?.skill3level}}级]
									</div>
								</div>
							</div>	
						</div>
						<div class="report-result">
							<img src="/src/assets/img/3d_report_txt_s_01.png" v-if="item.result == 2 || item.result == 1">
							<img src="/src/assets/img/3d_report_txt_s_02.png" v-if="item.result == 6">
							<img src="/src/assets/img/3d_report_txt_s_03.png" v-if="item.result == 0">
							<div>{{ formatTimestamp(item.time) }}</div>
						</div>
						<div class="report-hero-right">
							<div v-for="(hero,heroindex) in item.defend_all_hero_info">
								<div v-if="hero.id != 0">{{  herocfg[hero.id]?.name ? herocfg[hero.id]?.name : hero.id }}</div>
								<div class="report-hero" v-if="hero.id != 0">
									<div class="hero" style="height: 140px;">
										<img :src="`https://g0.gph.netease.com/ngsocial/community/stzb/cn/cards/cut/card_medium_${herocfg[hero.id]?.iconId ? herocfg[hero.id]?.iconId : hero.id }.jpg?gameid=g10`">
										<div class="bottom">
											<div>LV.{{ hero.level }}</div>
											<div>{{ hero.army }}</div>
										</div>
										<div class="stars">
											<div class="star" :class="{'up': hero.advance[0] >= (i+1)}" v-for="(s,i) in herocfg[hero.id]?.quality ? herocfg[hero.id]?.quality : 0"></div>
										</div>
									</div>
								</div>	
								<div class="skills" v-if="hero.id != 0 && item.all_skill_info.length > 0">

									<div class="skill" v-if="item.all_skill_info[getrightheroskillindex(heroindex)]?.skill1id != 0">
										{{ skillcfg[item.all_skill_info[getrightheroskillindex(heroindex)]?.skill1id]?.name ? skillcfg[item.all_skill_info[getrightheroskillindex(heroindex)]?.skill1id]?.name : '未知' }} 
										[{{ item.all_skill_info[getrightheroskillindex(heroindex)]?.skill1level}}级]
									</div>
									<div class="skill" v-if="item.all_skill_info[getrightheroskillindex(heroindex)]?.skill2id != 0">
										{{ skillcfg[item.all_skill_info[getrightheroskillindex(heroindex)]?.skill2id]?.name ? skillcfg[item.all_skill_info[getrightheroskillindex(heroindex)]?.skill2id]?.name : '未知' }} 
										[{{ item.all_skill_info[getrightheroskillindex(heroindex)]?.skill2level}}级]
									</div>
									<div class="skill" v-if="item.all_skill_info[getrightheroskillindex(heroindex)]?.skill3id != 0">
										{{ skillcfg[item.all_skill_info[getrightheroskillindex(heroindex)]?.skill3id]?.name ? skillcfg[item.all_skill_info[getrightheroskillindex(heroindex)]?.skill3id]?.name : '未知' }} 
										[{{ item.all_skill_info[getrightheroskillindex(heroindex)]?.skill3level}}级]
									</div>
								</div>
							</div>		
						</div>
					</div>
				</div>
			</div>

			<n-spin size="large" v-show="loaddataing && havemore" />
			<div v-show="havemore == false">没有更多数据了~</div>
		</div>
		<n-global-style />
  	</n-config-provider>
	
</template>

<style scoped>
	.report{
		width: 100%;
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.report-item{
		width: 100%;
		display: flex;
		background-color: rgba(0, 0, 0, .3);
		/* padding: 2rem; */
		box-sizing: border-box;
		flex-direction: column;
		margin: 8px 0;
		align-items: center;
	}

	.report-content{
		width: 100%;
		display: flex;
	}

	.report-hero{
		position: relative;
		margin: 0 8px;
	}

	.report-result{
		width: 15%;
		flex-shrink: 0;
	}

	.report-result>img{
		width: 150px;
		height: 150px;
		object-fit: cover;
	}

	.report-hero>.hero>img{
		width: 140px;
		height: 140px;
		object-fit: cover;
		object-position: left top;
	}


	.report-hero-right{
		flex-direction: row-reverse;
	}

	.report-hero-left,.report-hero-right{
		position: relative;
		display: flex;
		width: 42.5%;
		flex-shrink: 0;
		justify-content: flex-end;
		align-items: flex-start;
	}

	.report-hero>.hero>.bottom{
		width: 100%;
		display: flex;
		position: absolute;
		bottom: 0;
		left: 0;
		justify-content: space-between;
		padding: 0 2px;
		box-sizing: border-box;
		background-color: rgba(0, 0, 0, .6);
	}

	.report-header{
		display: flex;
		/* justify-content: center; */
		width: 100%;
	}

	.report-header>.left,.report-header>.right{
		width: 42.5%;
		flex-shrink: 0;
		/* flex: 1; */
		display: flex;
		justify-content: flex-end;
		flex-direction: column;
    	align-items: flex-end;
	}

	.report-header>.left>div,.report-header>.left>.info>div{
		margin-right: 8px;
		margin-left: 16px;
	}

	.report-header>.right>.info{
		display: flex;
		flex-direction: row-reverse;
	}

	.report-header>.left>.info{
		display: flex;
	}

	.report-header>.right>div,.report-header>.right>.info>div{
		margin-left: 8px;
		margin-right: 16px;
	}

	.report-header>.right{
		align-items: flex-start;
	}

	.report-header>.center{
		width: 15%;
		flex-shrink: 0;
		/* flex: 1; */
	}

	/* .report-header>div{
		margin: 0 16px;
	} */
	.report-bar{
		display: flex;
		position: relative;
		width: 100%;
		align-items: center;
	}

	.report-bar>.bar-box{
		display: flex;
		width: 70%;
		position: relative;
		height: 8px;
		background-color: rgba(0, 0, 0, .6);
	}

	.report-bar>.bar-box>.bar1{
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 8px;
		background-color: #633131;
	}

	.report-bar>.bar-box>.bar2.left,.report-bar>.bar-box>.bar1.left{
		right: 0;
    	left: unset;
	}

	.report-bar>.bar-box>.bar2.left{
		background-color: #4b84be;
	}

	.report-bar>.bar-box>.bar1.left{
		background-color: #204061;
	}

	.report-bar>.bar-box>.bar2{
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 8px;
		background-color: #a84141;
	}

	.stars {
		display: flex;
		position: absolute;
		top: 4px;
    	right: 8px;
	}

	.stars>.star{
		width: 20px;
		height: 20px;
		background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAMAAABg3Am1AAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAMAUExURQAAABoMBhoMBhoMBhoMBhkLBRoMBhoMBhoMBhoMBhoMBhUIAxoMBhoLBRkLBRcKBBYJBBQHAxkLBRUHBBoMBhoMBhoMBhcKBBgKBRkIBBoMBhoMBhgLBRkLBRMFAhoMBhkMBRMFAhQIBBAEARoMBhkMBhIEAhYGBCYUChYJBBkLBRgLBRgLBS0SCiEKBhoMBioWCyIRChoMBhoMBhcIA00lE2MxGCITChgLBTwlFFkqFXBcSI5KIzcZDUo3KHI7HScPCBgGBF9LOTkoGUUiET0jEz0bDrVrMkcyGoNDILFoL5+Qdm0+Gy0bD4d3ZpFbJnBVLL2pWXtfMHY9HE8yGFotFRMGBGFGI5VPJlo9IJl7PpF/aKleK1Y9JLuXSINjQZ1XKcy+p9bKq6STgr+wm6NoMXdmT4pRJn9LJMCyi825YDooF15MOj8eD4BtTrGiZU49LLCVUbGhfYVlMa+fjKqEP5tjLP/SYf/vbf//d//gZ//8c//pavrGW//zb///dP/jaP/NXv/tbP/5cf/ZZP/QYPfDWuy6V9ynTv/dZf//dv//fv//e//2cfXDXP//lv//+P//8//1bv/mafvJXv//hfW+WOeyU7hiLf/5Yv/ybP//6/7KXf//bfC8V+y1VK9fLP//1+GsUP///v/XYv//gv/xbf33hPvzg9CNQ/nQZPC3Vf//48ZtMcFqMP//s//+pd+FPP/8ZOuURP/+29+LQP//jv//ifK/WKNWKP/3e9Z8OezDW75xNNWHPr5kLcl3N6pZKeKwUf/8g9mBO+/pzP/5bOaiStqRQPTJZefCU+KZOvzMXf3YbtqLP/HPZNnEYMhvM97Tu8OCPc2ZR8uFPu2YRMeMQv7SZ///x+W6W/bu2+bbweuPQf74sP//xuWPQv341eygRu/bmu/NgPrubt+VPv//w/z6jP/gcPjqhN18OP/sef//zOqrT/rue+vgyv755v35n+vedrd0NuO1VemKP9mdSrNZKdqvbP76eOWbNPnwnvjvY96wQf/aaC/iuB8AAAB3dFJOUwAMBzUKVBsCAQMWbBE6P2B/ek+DHyExaGWhDilZV4ASQ4ZvdiNJjpagc0t5XL61K6yQJS2L0eSXSLXgz/jCs+evqb+lxsDK/cDz/unimtnx2fzi7MzYk8721u/f+sT97Pn8/On29tDt6Pb7rcvH3PC+9e/k7/TxZRQ62gAABRJJREFUSMfFVmVUHFcY7coMO7M2K7ACaxBcswsBQoh7iDbSuCd1b3dhfWGFwBoOiweCBCfu7u4ujdTd9fQtOW3TMKfZ/mnvOfPeOXPune+97/veffPMM/8TKEQCkfIv+EQ0MEAu9F3hBzOYPJGC7LOAQOINTMDYMOJzADY2fnEylS70UUAGASadG8OB/H0LgXAh8dCDbfMighk0nwRCOjV5kvv21zNlQXJfQiD+Is7rB923Nr8RESmV+JBaGiM44k133VdNH81NZCqITw8gD5K9fNfTWdy8aXoCKwp+agiJNDLiBfenmxo/aGqaFsYn/VP1KAhCpJGYiWl1bV2ZxZua97+YikFyAhHpv3UK4JIlsH88SRrDemmZ+/3GzMbizc37Z4RTIQYpkIsKgYry14cJKCxX0KVREJMawhG8uqHNqc/LyytubjozPkwcwothpzBI8f4wjej3qLAD6CnsID6VJRs4edTQkfNveDrzM4Egr/Hzba+ljRw6KjlRHEJlQlGMQBrS1zlRkRgHcEemTVm+4tx119q6zVV5+V5sc9Z9+cXSt+evAqrJAznRQQpvd5Hj+RELV01ZvvTGdffatS63p85ZtSazD3rnVY/b5fJsOLhsxaQpaanRUbC3lwN4wxffq/O4XO62q51dzmJ9lT5/DQAQ6POLnV2dtzd4XO67t36ci0Fc7+mSi1hxUx3373U5fzpUXl7eU6UGVLVarwcDeNFzyL7f+cv9o4czBJFS1JskQkAMFjf15+PHenrKKyoKstWAnG3P7oO9oEJlqqj69bfjDzME0ZC8r4x+QhITi3v3vSPH1GqTqkC3DmjsuoICnU5n16lMGpP6ge3h8edHsKBA8qNa+IHyimNXehV2lV1VaNLpTH1QqUwmTWH2A9vGXsAPiv+zTYCCz4ldeenIjmz7Oo1Bo9IYDBpNYaF3Klynqt94Z/Vwlij+sbaiSOh8TvrES6d26PrIhqwsA3i8MJnqu3tXx2ExCsLfzj3K4MteGdd7qlpTmJWblZubm9OQ4x2zDPsO3JwK+AFPOAIFZfASR4/7+FSroSE3R6lsMJuVDQC56w90L4kVM0n9bI2CSqlhMy+XtDYYlUqj2WwEk9KsVB7uXpDOYZJoFDyzCEmdvu2KsVJrOWE2ak8YjZWVlUWHv5sgCE5B8Q6eJBRbeOaTIqO5yOGwaLXaIq3WUu347OZzw4MZBLwDh7Jlc05/6NCajRZr6ZbKIotli7VUu763dkxICg3fv8KmnTxpLa2pKa23FimV3gBWZcmdstFiNt6SgF2Ez2g9ai212UqrrfW2khyDsqTaXLJxd5IMgv1w7gQFM+KtHdeqS6zfWmtadu7ZVX+61VFtqdl1dlaYiIvgWnDCPMc1h3bfhYvdBzrK2s/uabHVbNl3sT0jPAbPMwmM4IR3rhx1fP/NxN0dtcNmPbu1rGPn+Qu2lvZhAlwHFIayUhe12nZ11O4dMmisYPaIJCApm3i+ZfuEBFw/k0SJxy+yte8G9HQBRuUFy0aMHbZkb+3O7fiFoMBszpzLe34YMjg2nEWF6AGMIKpYMHbQhK3bQSFC+xcCGSCanTHOS4/mQXSukEyTAwkmiB28oCyJw0bxBMlJg+PCIvlsEkxAvFe1ZABdRI0OTx80mhMq6b8klEHFsEhmaABM/iOHCM2fDvGiMRbuphGYzoakCpT8eMoRIZcUCrFJqB9eM0m4XJT8ZNMgQtgfliD/4W/L76xOxS1MHlMIAAAAAElFTkSuQmCC) no-repeat;
		background-size: 100% 100%;
	}

	.stars>.star.up{
		background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAMAAABg3Am1AAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAMAUExURQAAABkLBRoMBhoMBhoMBhoMBhkLBRoMBhoMBhoMBhoMBhoMBhoMBhcLBRULBRoMBhoLBRMJBRMLBRoMBhgLBhQJBRYIBBoMBhcLBRgMBhoMBhMJBBYKBRkIBBkLBRMKBBUKBb8PAxYLBRoMBi0IAx4HAxIKBUUIBBQKBBgMBlsLA3gPBCILBiEIBDUNBs4UBA4LBpkTBIYJAykOCVIIA40QBBoMBmQPBG4RBCYDADYJAkcXB0wGAl8hDDoIBJ8NA7ENA94UBNMvDGIuHksOBEIYCKhRH9+fSkAJA7ddKtWAPqYxDKwRBE4HBOUWBYUUBUIqIHBEM6BrUc2kfVMlE+JJEmAKA3AvE4dFIHotDX8jDZMvDsc2DZBCG7svC8YiCLchCN4kCKkhCH4UBbJyS4ddStuue8yXY7xxOsZ1NpM3GdWRMMF2KdehWpU2FdZBEa02DZhKG2keDZtDHcxXF79GEf5dGfxMFP9xHv94IP+LJPxJE/97IPdFEvtUFv+IJP+EJP9iGv9XF/9sHf+AIv+XKP+pLf91H/s7EP9qHPMvDP+hLP+2Mv5UFvY3Dv5IE/9mG/c7D//YPfY/EPlPFe8YBfcrC/+TJv+cKv/NOv7sff1XF/9RFfpDEf9OFPI1DfYdB/UxDOsWBf5AEf9+IfdiGf/6dv+zLf/SOv/FOP+sKP/KZP+9NP/4Tf+8ZvEfB+0jCfQmCewcB+xAEf6QJv+kLf//vP+xMv+qOP//pvdzH/7tWPKDP/6hJf/ab/BoG+8sC/AmCvI8D/MhCeUyDPRSFv+vLf7zav+mQP/3gv/+XP/+Z/+sWfR8H/JcFv/BRP/CMN8wDPMVBerCb//+kPxiG+5ZFvnnkPNxG//0Xv/fdP/gQP/Njf/sR/+xYv/RSPhzKP+mS/6fUv+zTv//l/+4UvqIL//AOf8dBvrqsfntqv7saPWDKO/BUf7ybPK5Wuy6O//jaf+TMv/NdOOMSv+RO/+STf6ePP/wkOuqUvBaGfKrT/iUO+akSv/lTu1fGS7CN6cAAAB2dFJOUwBDDBIPGTgDAQUVKAhTYyE1a18dWXyCAk9LLG1/djxnVulcJIyEep5wRrTFg4mZ8XLWyZCp0zC9xJKeqbu1ldni+Pi1rqTg/a7q+e3fpfzNm7rZ763+wcDNxsrX99vz7uf63s7mxPny7PLa+u332vryy8Da8Op7pU6XAAAFLUlEQVRIx8VWVVRbWRRt3F2IGwR3aCkF6u7u7UjHXeNKhJAEAkmQQCF4cddCoQZ1m6m7d9xd1poXumamq2RNMz8z7+N+7X33O3efs++dNOl/+hBgKBjxL/AoNIvIggTOQMBAeJJoCjhgAhROWrkqWhCMClSAgIt+fs3KcBAaGbBA0pqPV8eQMagABQTiVZ/0vrCWGwELCkiAx459sXfgwuokPhEbiECwQLzxQu/Zsx/O4TAoT5ZAQkDh817vHThSV/JcEh/+5KPFssgxG0sGjhz6tA6QwBGe5B4SxgxNW3C3bnDwUN0XL8WSeNB/AiNQWAgcn7qtZOBYy9G6Q5c+3yCmsqBY1MRCkMggAAsLZhFBZNrcBXePtLQcPXbs0tHXprEFICKGgIZiEYi/XUSB0ZRgFhzEFJBJXPq0+SU/DTrrr7W0XBv8alMqncsX4ZggOCsYBsU+lMIS4BEMKp7NFadOS8t4592fS4adp5sbPE6n58qr2zasnRcbEs0lkQUMHgbqOwIEhcGNjkmKTcvY/Mqixd8N3W5tunHa3tBgb3B6bnzd9OWCl+dv2ZQxzgIKAgjgKfy0qZvfWLR4aOj61avfHv/h1qnTY+YGZ73zSv2ppuHLl5uO3/7198WL5r89J5RJedhqC5fc/P5W6/Xjra13Dn9TvKu52W62ezyeek+zva//8OF7TcMlw6337m8BbPGVgKHS0t/Kvjn0Y3+/uXFP0Z5Gs73PZjL3mcfsfY27xnY1P3hw6s5v93cvWcgFwXw1QIh4Tvr2kycv7mm09fQUF5u6i01ao8KkVdjG+uxmmxnY5sTukSULaQLMeDMi0HA8PX1r7ckTPTaTqdikMBqNWVm6boWiqKffptApTN0jI8b359KoLPBDL8YZ69/MO3/C1m3UKrQ6lVwmU+m0OrmxSFGjKC4a+cXwXgZNxPqrTYLQPD596rpz5y/WyOUqralIZZDJ5TJg0dXItLbPake3R3LIj4YCAmDERK1zd40CIJWuSJcF4MdVjDW63bXntqZzyETooxMehAbxY3Yc3H9wNDsrS6ZSAUu2vEZl1OkAvGZ2Cv0xvC+KIkhJO87s75LIDHq9Pi9Pn2fIzs6Sq2pHNbPX0/FwyOMJgqJEsFMnd2kkkszMnEyJRKLXS/Jksmy168xUOp6HntjiqGARLaUrD4BrNKV6SU6O3mAw5Ond+XFz45kwpL9Zw0WvUGc6yjNL3RpJqVRaKinVaCSasmci40H+5g5BwYUkOMqV0pydSrdb6pBKlWVqN8D0RtGYEH+JjRHFPutwqNVql0vtlpbnlDvKlICCOneFmOHvl1AsfNjT1dVq136vy+1Wu5Q7c8ql7kxlbmVCiN/wwBL5yU9Z8vMt+a6yMu/e9n250p1SqcOSW7l0OtVfkIN5pMjZ1fkWS67VWtjRVlDVvjdXrbR4rYXCMDLLDwEKik85WG2xWNs72yoKhEvjCjr27fV6vdaOuGS8v5CFMEOjrNWV+zqrCqqEicmxkTPi2ira2q2VhbMiSTywHxsY4smFlZ0FVcuFicvo4aT4kMiEuIKKAx2Fy1P8GQHcCiEJBw58MGtmVBiNTQXBI0Th4uTJwqqKwoqoUD9GoDDkMKEPPj2UROURIFgIBiQKpy9LFHZ+NCMGB/NLSJjhg+PgBF9gIVEAhcqmhSXOTAzx4xyCEhHO4XDxDCLlzysdoBB4VHYoh4YngicQkFgCD4cDTYGBH21kBJQAZ+IYcLQfp5FYNIGAnvBgQEFhBAIE8R8+W/4AABnMxy2KV2sAAAAASUVORK5CYII=);
	}
</style>
