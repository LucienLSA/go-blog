<template>
  <div class="valid-email-result">
    <div v-if="loading">正在验证，请稍候...</div>
    <div v-else>
      <div v-if="success" style="color:green;">邮箱绑定成功！</div>
      <div v-else style="color:red;">邮箱绑定失败：{{ msg }}</div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loading: true,
      success: false,
      msg: ''
    }
  },
  mounted() {
    const token = this.$route.params.token;
    this.$axios({
      method: 'get',
      url: `/api/v2/user/validEmail?token=${token}`
    }).then(res => {
      this.loading = false;
      if (res.code === 1000) {
        this.success = true;
        this.msg = '邮箱绑定成功！';
      } else {
        this.success = false;
        this.msg = res.msg || '邮箱绑定失败';
      }
    }).catch(() => {
      this.loading = false;
      this.success = false;
      this.msg = '网络错误';
    });
  }
}
</script>

<style scoped>
.valid-email-result {
  max-width: 400px;
  margin: 80px auto;
  padding: 40px 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.08);
  text-align: center;
  font-size: 18px;
}
</style> 