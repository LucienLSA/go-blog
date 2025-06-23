<template>
  <div class="user-info-card">
    <div class="avatar-section">
      <img :src="user.avatar || require('@/assets/images/avatar.png')" alt="avatar" class="avatar-img" @click="triggerAvatarUpload" />
      <input type="file" ref="avatarInput" class="avatar-upload" @change="uploadAvatar" accept="image/*" />
      <div class="avatar-tip">点击头像更换</div>
    </div>
    <form @submit.prevent="updateUser" class="info-form">
      <div class="form-group">
        <label>用户名</label>
        <input v-model="user.user_name" />
      </div>
      <div class="form-group">
        <label>年龄</label>
        <input v-model.number="user.age" type="number" min="1" max="130" />
      </div>
      <div class="form-group">
        <label>性别</label>
        <select v-model="user.gender">
          <option value="男">男</option>
          <option value="女">女</option>
          <option value="未知">未知</option>
        </select>
      </div>
      <div class="form-group">
        <label>邮箱</label>
        <input v-model="user.email" type="email" placeholder="输入邮箱地址" class="full-input" :disabled="isEmailVerifying" />
      </div>
      <div class="form-group btn-group" v-if="!isEmailVerifying">
        <label></label>
        <button type="button" class="secondary" @click="handleSendEmail(1)" :disabled="!user.email">绑定</button>
        <button type="button" class="secondary" @click="handleSendEmail(2)" :disabled="!user.email">解绑</button>
      </div>
      <div v-if="isEmailVerifying">
        <div class="form-group">
          <label></label>
          <input v-model="emailCode" maxlength="6" placeholder="请输入验证码" class="full-input" />
        </div>
        <div class="form-group btn-group">
          <label></label>
          <button type="button" class="primary" @click="submitEmailCode">确认</button>
          <button type="button" class="secondary" @click="cancelEmailVerify">取消</button>
        </div>
      </div>
      <div class="form-group">
        <label>密码修改</label>
        <button type="button" class="secondary" @click="showPasswordForm = !showPasswordForm">
          {{ showPasswordForm ? '取消修改' : '修改密码' }}
        </button>
      </div>
      <div v-if="showPasswordForm">
        <div class="form-group">
          <label>当前密码</label>
          <input v-model="passwords.current" type="password" class="full-input" placeholder="请输入当前密码" />
        </div>
        <div class="form-group">
          <label>新密码</label>
          <input v-model="passwords.new" type="password" class="full-input" placeholder="请输入新密码" />
        </div>
        <div class="form-group">
          <label>确认密码</label>
          <input v-model="passwords.confirm" type="password" class="full-input" placeholder="请再次输入新密码" />
        </div>
      </div>
      <div class="form-actions">
        <button type="submit" class="primary">保存修改</button>
      </div>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      user: {
        user_name: "",
        age: 18,
        gender: "男",
        email: "",
        avatar: ""
      },
      currentOpType: 1,
      emailCode: '',
      isEmailVerifying: false,
      showPasswordForm: false,
      passwords: {
        current: '',
        new: '',
        confirm: ''
      }
    }
  },
  created() {
    this.loadUserInfo();
  },
  methods: {
    loadUserInfo() {
      const loginResult = JSON.parse(localStorage.getItem("loginResult"));
      if (loginResult) {
        this.user.user_name = loginResult.user_name;
        // TODO: 加载用户完整信息的API调用
      }
    },
    updateUser() {
      // 构建更新数据
      const updateData = {};
      
      // 只传递已修改的字段
      if (this.user.user_name) updateData.user_name = this.user.user_name;
      if (this.user.age) updateData.age = this.user.age;
      if (this.user.gender) updateData.gender = this.user.gender;
      if (this.user.email) updateData.email = this.user.email;

      // 如果正在修改密码，添加密码相关字段
      if (this.showPasswordForm) {
        if (!this.passwords.current || !this.passwords.new || !this.passwords.confirm) {
          alert('请填写完整的密码信息');
          return;
        }
        if (this.passwords.new !== this.passwords.confirm) {
          alert('两次输入的新密码不一致');
          return;
        }
        updateData.password = this.passwords.current;
        updateData.new_password = this.passwords.new;
      }

      this.$axios({
        method: "put",
        url: "/api/v2/user/update",
        data: updateData
      }).then(res => {
        if (res.code === 1000) {
          alert("修改成功");
          this.showPasswordForm = false;
          this.passwords = { current: '', new: '', confirm: '' };
        } else {
          alert(res.msg || "修改失败");
        }
      }).catch(err => {
        alert(err.response?.data?.msg || "修改失败");
      });
    },
    uploadAvatar(e) {
      const file = e.target.files[0];
      if (!file) return;
      
      const formData = new FormData();
      formData.append("file", file);
      
      this.$axios({
        method: "post",
        url: "/api/v2/user/upload/avatar",
        data: formData,
        headers: { "Content-Type": "multipart/form-data" }
      }).then(res => {
        if (res.code === 1000) {
          this.user.avatar = res.data.avatar;
          alert("头像上传成功");
        } else {
          alert(res.msg || "头像上传失败");
        }
      }).catch(err => {
        alert(err.response?.data?.msg || "头像上传失败");
      });
    },
    handleSendEmail(type) {
      if (!this.user.email) {
        alert('请输入邮箱地址！');
        return;
      }
      this.currentOpType = type;
      this.isEmailVerifying = true;
      
      this.$axios({
        method: 'post',
        url: '/api/v2/user/sendEmail',
        data: {
          operation_type: type,
          email: this.user.email
        }
      }).then(res => {
        if (res.code === 1000) {
          alert("验证码已发送，请查收邮箱");
        } else {
          alert(res.msg || '验证码发送失败，请稍后重试');
          this.isEmailVerifying = false;
        }
      }).catch(err => {
        alert(err.response?.data?.msg || '验证码发送失败');
        this.isEmailVerifying = false;
      });
    },
    submitEmailCode() {
      if (!this.emailCode || this.emailCode.length !== 6) {
        alert('请输入6位验证码！');
        return;
      }
      
      this.$axios({
        method: 'post',
        url: '/api/v2/user/validEmail',
        data: {
          user_name: this.user.user_name,
          email: this.user.email,
          code: this.emailCode,
          operation_type: this.currentOpType
        }
      }).then(res => {
        if (res.code === 1000) {
          alert(this.currentOpType === 1 ? '邮箱绑定成功！' : '邮箱解绑成功！');
          this.emailCode = '';
          this.isEmailVerifying = false;
          if (this.currentOpType === 2) {
            this.user.email = '';
          }
        } else {
          alert(res.msg || '验证码错误');
        }
      }).catch(err => {
        alert(err.response?.data?.msg || '验证失败，请重试');
      });
    },
    cancelEmailVerify() {
      this.isEmailVerifying = false;
      this.emailCode = '';
    },
    triggerAvatarUpload() {
      this.$refs.avatarInput.click();
    }
  }
}
</script>

<style scoped>
.user-info-card {
  max-width: 400px;
  margin: 48px auto;
  background: #fff;
  border-radius: 18px;
  box-shadow: 0 4px 24px rgba(0,0,0,0.07);
  padding: 36px 32px 28px 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 18px;
}
.avatar-img {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  border: 2px solid #f0f0f0;
  object-fit: cover;
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.avatar-img:hover {
  box-shadow: 0 0 0 4px #e6f7ff;
}
.avatar-upload {
  display: none;
}
.avatar-tip {
  font-size: 13px;
  color: #aaa;
  margin-top: 6px;
}
.info-form {
  width: 100%;
}
.form-group {
  display: flex;
  align-items: center;
  margin-bottom: 18px;
}
.form-group label {
  width: 80px;
  text-align: right;
  margin-right: 16px;
  font-size: 16px;
  color: #222;
  flex-shrink: 0;
}
.full-input {
  width: 320px;
  height: 38px;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
  padding: 0 12px;
  font-size: 15px;
}
.full-input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}
.btn-group {
  align-items: flex-start;
  gap: 12px;
}
.secondary, .primary {
  height: 38px;
  border-radius: 6px;
  font-size: 15px;
  cursor: pointer;
  min-width: 80px;
  margin-right: 12px;
}
.secondary {
  border: 1px solid #4caf50;
  background: #fff;
  color: #4caf50;
  padding: 0 18px;
  transition: background 0.2s;
}
.secondary:hover:not(:disabled) {
  background: #e8f5e9;
}
.secondary:disabled {
  border-color: #ccc;
  color: #999;
  cursor: not-allowed;
}
.primary {
  border: none;
  background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
  color: #fff;
  padding: 0 24px;
  transition: opacity 0.2s;
}
.primary:hover:not(:disabled) {
  opacity: 0.9;
}
.primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.form-actions {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style> 