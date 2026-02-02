/*
 * gin-vue-admin web框架组
 *
 * */
// 加载网站配置文件夹
import { register } from './global'
import packageInfo from '../../package.json'

export default {
  install: (app) => {
    register(app)
    console.log(`
       欢迎使用 LoRAForge
       当前版本:v${packageInfo.version}
       ** 感谢您对LoRAForge的支持与关注 合法授权使用更有利于项目的长久发展**
    `)
  }
}
