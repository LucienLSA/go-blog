<template>
  <div class="main">
    <div class="container">
      <h2 class="form-title">登录</h2>
      <div class="form-group" v-if="loginType === 'password'">
        <label for="name">用户名</label>
        <input type="text" class="form-control" v-model="username" name="name" id="name" placeholder="用户名" />
      </div>
      <div class="form-group" v-if="loginType === 'password'">
        <label for="pass">密码</label>
        <input type="password" class="form-control" v-model="password"  name="pass" id="pass" placeholder="密码" />
      </div>
      <div v-if="loginType === 'email'">
        <div class="form-group">
          <label for="email">邮箱</label>
          <input type="email" class="form-control" v-model="email" name="email" id="email" placeholder="邮箱" />
        </div>
        <div class="form-group">
          <label for="code">验证码</label>
          <div style="display:flex;align-items:center;">
            <input type="text" class="form-control" v-model="emailCode" name="code" id="code" placeholder="验证码" style="flex:1;" />
            <button type="button" class="btn btn-info" style="margin-left:10px;" @click="sendEmailCode" :disabled="codeBtnDisabled">{{codeBtnText}}</button>
          </div>
        </div>
      </div>
      <div class="form-btn">
        <button type="button" class="btn btn-info" @click="submit">提交</button>
      </div>
      <div style="text-align:center;margin-top:10px;">
        <a href="javascript:void(0);" @click="toggleLoginType">
          {{ loginType === 'password' ? '使用邮箱验证码登录' : '使用账号密码登录' }}
        </a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
	name: "Login",
	data() {
		return {
			username: "",
			password: "",
			submitted: false,
			loginType: 'password',
			email: '',
			emailCode: '',
			codeBtnText: '获取验证码',
			codeBtnDisabled: false,
			codeTimer: null,
			codeCountdown: 60
		};
	},
	computed: {
	},
	created() {

	},
	methods: {
		toggleLoginType() {
			this.loginType = this.loginType === 'password' ? 'email' : 'password';
		},
		sendEmailCode() {
			if (!this.email) {
				alert('请输入邮箱');
				return;
			}
			this.codeBtnDisabled = true;
			this.codeBtnText = '发送中...';
			this.$axios({
				method: 'post',
				url: '/api/v2/user/login/emailcode',
				data: JSON.stringify({
					email: this.email,
					operation_type: 1
				})
			}).then((res) => {
				if (res.code == 1000) {
					this.startCodeCountdown();
				} else {
					alert(res.msg || '发送失败');
					this.codeBtnDisabled = false;
					this.codeBtnText = '获取验证码';
				}
			}).catch(() => {
				this.codeBtnDisabled = false;
				this.codeBtnText = '获取验证码';
			});
		},
		startCodeCountdown() {
			this.codeCountdown = 60;
			this.codeBtnText = this.codeCountdown + 's';
			this.codeTimer = setInterval(() => {
				this.codeCountdown--;
				if (this.codeCountdown <= 0) {
					clearInterval(this.codeTimer);
					this.codeBtnText = '获取验证码';
					this.codeBtnDisabled = false;
				} else {
					this.codeBtnText = this.codeCountdown + 's';
				}
			}, 1000);
		},
		submit() {
			if (this.loginType === 'password') {
				this.$axios({
					method: 'post',
					url:'/api/v2/user/login',
					data: JSON.stringify({
						user_name: this.username,
						password: this.password
					})
				}).then((res)=>{
					if (res.code == 1000) {
						localStorage.setItem("loginResult", JSON.stringify(res.data));
						this.$store.commit("login", res.data);
						this.$router.push({path: this.redirect || '/' })
					} else {
						alert(res.msg || '登录失败');
					}
				}).catch(()=>{});
			} else {
				this.$axios({
					method: 'post',
					url:'/api/v2/user/login/email',
					data: JSON.stringify({
						email: this.email,
						code: this.emailCode
					})
				}).then((res)=>{
					if (res.code == 1000) {
						localStorage.setItem("loginResult", JSON.stringify(res.data));
						this.$store.commit("login", res.data);
						this.$router.push({path: this.redirect || '/' })
					} else {
						alert(res.msg || '登录失败');
					}
				}).catch(()=>{});
			}
		}
	}
};
</script>
<style lang="less" scoped>
.main {
  background: #f8f8f8;
  padding: 150px 0;
  .container {
    width: 600px;
    background: #fff;
    margin: 0 auto;
    max-width: 1200px;
    padding: 20px;
    .form-title {
      margin-bottom: 33px;
      text-align: center;
    }
    .form-group {
      margin: 15px;
      label {
        display: inline-block;
        max-width: 100%;
        margin-bottom: 5px;
        font-weight: 700;
      }
      .form-control {
        display: block;
        width: 100%;
        height: 34px;
        padding: 6px 12px;
        font-size: 14px;
        line-height: 1.42857143;
        color: #555;
        background-color: #fff;
        background-image: none;
        border: 1px solid #ccc;
        border-radius: 4px;
      }
    }
    .form-btn {
      display: flex;
      justify-content: center;
      .btn {
        padding: 6px 20px;
        font-size: 18px;
        line-height: 1.3333333;
        border-radius: 6px;
        display: inline-block;
        margin-bottom: 0;
        font-weight: 400;
        text-align: center;
        white-space: nowrap;
        vertical-align: middle;
        -ms-touch-action: manipulation;
        touch-action: manipulation;
        cursor: pointer;
        border: 1px solid transparent;
      }
      .btn-info {
        color: #fff;
        background-color: #5bc0de;
        border-color: #46b8da;
      }
    }
  }
}
</style>