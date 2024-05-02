import { GetConfigs, GetENVConfigs, InitConfig, SaveENVConfigs } from 'wailsjs/go/main/GOContext'
import type { common } from 'wailsjs/go/models'
import { WindowSetTitle } from 'wailsjs/runtime/runtime'

function createDataStore() {
  const [envConfig, setEnvConfig] = createStore<common.YamlInfo>({} as common.YamlInfo)
  const [configs, setConfigs] = createStore<common.TypeConfig[]>([])
  async function getEnv() {
    GetENVConfigs().then((res) => {
      setEnvConfig(res)
    })
  }
  async function updateEnvConfig(name: any, value: string) {
    setEnvConfig(name, value)
  }
  async function updateEnv(key: string, index: number) {
    setEnvConfig('env', (env) => {
      env[key].current = index
      return env
    })
  }
  async function saveEnv() {
    if (!envConfig)
      return
    SaveENVConfigs(envConfig)
  }
  async function getData() {
    GetConfigs().then((result) => {
      if (!result || !result.length)
        return
      // 工具类别按照index进行排序
      result.sort((a, b) => {
        if (a.index === 0)
          a.index = 1e10
        if (b.index === 0)
          b.index = 1e10
        return a.index - b.index
      })
      // 类别里的工具按照index进行排序
      result = result.map((c) => {
        c.config.sort((a, b) => {
          if (a.index === 0)
            a.index = 1e10
          if (b.index === 0)
            b.index = 1e10
          return a.index - b.index
        })
        return c
      })
      setConfigs(result)
    })
  }
  async function refresConfig() {
    await InitConfig().catch((e: string) => {
      console.error(e)
      throw new Error(e)
    })
    await getData()
    await getEnv()
    // 设置窗口标题
    if (envConfig.title)
      WindowSetTitle(envConfig.title)
  }
  return { configs, envConfig, refresConfig, saveEnv, updateEnv, updateEnvConfig }
}

export default createRoot(createDataStore)
