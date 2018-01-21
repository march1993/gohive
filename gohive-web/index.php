<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Gohive Control Gui</title>
	<link rel="stylesheet" href="css/index.css">
	<style>
		#app{
			width:980px;
			margin:0 auto;
		}
		.model{
			margin-top:10px;
		}
	</style>
</head>
<body>
	<div id="app">
		Enter your token: <el-input v-model="token" style="width:200px;"></el-input>
		<div class="model">
			<el-input v-model="app_name" style="width:200px;"></el-input>
			<el-button @click="CreateAPP()">Create APP</el-button>
			<el-alert
			    title=""
			    type="error"
			    :closable="false"
			    style="width:400px;" v-if="create_flag">{{create_warning}}
			</el-alert>
			<el-alert
			    title=""
			    type="success"
			    :closable="false"
			    style="width:200px;" v-if="create_success">Created Successfully
			</el-alert>
		</div>
		<div class="model">
			<el-input v-model="repair" style="width:200px;"></el-input><el-button @click="RepairAPP()" style="margin-left: 5px;">Repair APP</el-button>
			<el-alert
			    title=""
			    type="error"
			    :closable="false"
			    style="width:400px;" v-if="repair_flag">{{repair_warning}}
			</el-alert>
			<el-alert
			    title=""
			    type="success"
			    :closable="false"
			    style="width:200px;" v-if="repair_success">Repaired Successfully
			</el-alert>
		</div>
		<div class="model">
			<el-input v-model="status_app" style="width:200px;"></el-input><el-button @click="StatusAPP()" style="margin-left: 5px;">Status APP</el-button>
			<el-alert
			    title=""
			    type="error"
			    :closable="false"
			    style="width:400px;" v-if="status_flag">{{status_warning}}
			</el-alert>
			<ul>
				<li v-for="(item, index) in status_list" style="margin-bottom: 10px;">
					{{index}} : {{item.Status}} - {{item.Reason}}
				</li>
			</ul>
			
		</div>
		<div class="model">
			<el-button @click="GetAPPList()">Get APP List</el-button>
			<el-alert
			    title=""
			    type="error"
			    :closable="false"
			    style="width:400px;" v-if="list_flag">{{list_warning}}
			</el-alert>
			<ul>
				<li v-for="(item, index) in list" style="margin-bottom: 10px;">
					{{item}}
					<el-button @click="RemoveAPP(list_noprefix[index])">Delete</el-button>
					<el-input v-model="rename[index]" style="width:100px;"></el-input>
					<el-button @click="RenameAPP(index)">Rename</el-button>
				</li>
			</ul>
		</div>
	</div>
</body>
<script src="js/vue.js"></script>
<script src="js/index.js"></script>
<script src="js/axios.min.js"></script>
<script>
	new Vue({
		el: '#app',
		data: function() {
			return {
				token: "",
				app_name: "",
				create_flag: false,
				create_warning: "",
				create_success: false,
				list_flag: false,
				list_warning: "",
				list: "",
				list_noprefix: "",
				prefix: "",
				rename: [],
				repair: "",
				repair_flag: false,
				repair_success: false,
				repair_warning: "",
				status_app: "",
				status_flag: false,
				status_warning: "",
				status_list: ""
			}
		},
		methods: {
			CreateAPP: function () {
				if(this.app_name == "" || this.token == ""){
					this.create_warning = "APP名称或者token不能为空";
					this.create_flag = true;
				}
				axios.post('http://vultr1.xuxuxu.me/app/createApp', {
				    Token: this.token,
				    App: this.app_name
				})
				.then((response) => {
				    if(response.data.Status == "STATUS_SUCCESS"){
				    	this.create_flag = false;
				    	this.create_success = true;
				    }
				})
				.catch((error) => {
				    console.log(error);
				});
			},
			GetAPPList: function () {
				if(this.token == ""){
					this.list_warning = "token cannot be empty";
					this.list_flag = true;
				}
				axios.post('http://vultr1.xuxuxu.me/app/getAppList', {
				    Token: this.token
				})
				.then((response) => {
				    if(response.data.Status == "STATUS_SUCCESS"){
				    	this.list = response.data.Result;
				    	this.prefix = response.data.Addition.PREFIX;
				    	this.list_noprefix = response.data.Result.map(function (name) {return name.split(response.data.Addition.PREFIX)[1]});
				    }
				})
				.catch((error) => {
				    console.log(error);
				});
			},
			RemoveAPP: function (app) {
				axios.post('http://vultr1.xuxuxu.me/app/removeApp', {
				    Token: this.token,
				    App: app
				})
				.then((response) => {
				    if(response.data.Status == "STATUS_SUCCESS"){
				    	this.GetAPPList();
				    }
				})
				.catch((error) => {
				    console.log(error);
				});
			},
			RenameAPP: function (index) {
				//console.log(this.rename[index]);
				if(this.rename[index] == undefined || this.rename[index] == ""){
					this.list_flag = true;
					this.list_warning = "重命名不能为空";
					setTimeout((function(){this.list_flag = false;}).bind(this),2000);
				}
				axios.post('http://vultr1.xuxuxu.me/app/renameApp', {
				    Token: this.token,
				    OldName: this.list_noprefix[index],
				    NewName: this.rename[index]
				})
				.then((response) => {
				    if(response.data.Status == "STATUS_SUCCESS"){
				    	this.GetAPPList();
				    }
				})
				.catch((error) => {
				    console.log(error);
				});
			},
			RepairAPP: function () {
				if(this.repair == ""){
					this.repair_warning = "app to be repaired cannot be empty";
					this.repair_flag = true;
					setTimeout((function(){this.repair_flag = false;}).bind(this),2000);
				}
				axios.post('http://vultr1.xuxuxu.me/app/repairApp', {
				    Token: this.token,
				    App: this.repair,
				})
				.then((response) => {
				    if(response.data.Status == "STATUS_SUCCESS"){
				    	this.GetAPPList();
				    	this.repair = "";
				    }
				})
				.catch((error) => {
				    console.log(error);
				});
			},
			StatusAPP: function () {
				if(this.status_app == undefined || this.status_app == undefined){
					this.status_warning = "app to be checked status cannot be empty";
					this.status_flag = true;
					setTimeout((function(){this.status_flag = false;}).bind(this),2000);
				}else{
					axios.post('http://vultr1.xuxuxu.me/app/statusApp', {
					    Token: this.token,
					    App: this.status_app,
					})
					.then((response) => {
					    if(response.data.Status == "STATUS_SUCCESS"){
					    	this.status_list = response.data.Result;
					    }
					})
					.catch((error) => {
					    console.log(error);
					});
				}
			}
		}
	})
</script>
</html>