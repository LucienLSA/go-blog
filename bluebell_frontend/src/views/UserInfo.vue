<template>
  <div class="user-info">
    <h2>用户信息</h2>
    <form @submit.prevent="updateUser">
      <div>
        <label>用户名：</label>
        <input v-model="user.user_name" disabled />
      </div>
      <div>
        <label>年龄：</label>
        <input v-model="user.age" type="number" min="1" max="130" />
      </div>
      <div>
        <label>性别：</label>
        <select v-model="user.gender">
          <option value="男">男</option>
          <option value="女">女</option>
          <option value="未知">未知</option>
        </select>
      </div>
      <div>
        <label>邮箱：</label>
        <input v-model="user.email" />
        <button type="button" @click="sendBindEmail">绑定/解绑邮箱</button>
      </div>
      <div>
        <label>头像：</label>
        <img :src="user.avatar || require('@/assets/images/avatar.png')" alt="avatar" width="80" />
        <input type="file" @change="uploadAvatar" />
      </div>
      <button type="submit">保存修改</button>
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
        avatar: require('@/assets/images/user.png')
      }
    }
  },
  created() {
    // 假设登录后用户信息已存储在本地
    const loginResult = JSON.parse(localStorage.getItem("loginResult"));
    if (loginResult) {
      this.user.user_name = loginResult.user_name;
      // 其它信息可从后端拉取或本地存储
    }
  },
  methods: {
    updateUser() {
      this.$axios({
        method: "post",
        url: "/user/update",
        data: JSON.stringify({
          user_name: this.user.user_name,
          age: this.user.age,
          gender: this.user.gender,
          email: this.user.email
        })
      }).then(res => {
        if (res.code == 1000) {
          alert("修改成功");
        } else {
          alert(res.msg || "修改失败");
        }
      });
    },
    uploadAvatar(e) {
      const file = e.target.files[0];
      const formData = new FormData();
      formData.append("avatar", file);
      this.$axios({
        method: "post",
        url: "/user/avatar",
        data: formData,
        headers: { "Content-Type": "multipart/form-data" }
      }).then(res => {
        if (res.code == 1000) {
          this.user.avatar = res.data.avatar;
          alert("头像上传成功");
        } else {
          alert(res.msg || "头像上传失败");
        }
      });
    },
    sendBindEmail() {
      this.$axios({
        method: "post",
        url: "/user/sendEmail",
        data: JSON.stringify({
          operation_type: 1, // 1为绑定，2为解绑
          email: this.user.email
        })
      }).then(res => {
        if (res.code == 1000) {
          alert("邮件已发送，请查收邮箱完成绑定/解绑操作");
        } else {
          alert(res.msg || "操作失败");
        }
      });
    }
  }
}
</script>

<style scoped>
.user-info {
  max-width: 500px;
  margin: 40px auto;
  background: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}
.user-info form > div {
  margin-bottom: 18px;
}
.user-info label {
  display: inline-block;
  width: 80px;
  font-weight: bold;
}
.user-info input[type="text"],
.user-info input[type="number"],
.user-info input[type="email"],
.user-info select {
  width: 200px;
  padding: 4px 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.user-info button {
  margin-left: 10px;
  padding: 4px 12px;
  border: none;
  background: #54b351;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
}
.user-info img {
  vertical-align: middle;
  margin-right: 10px;
  border-radius: 50%;
  border: 1px solid #eee;
}
</style> 