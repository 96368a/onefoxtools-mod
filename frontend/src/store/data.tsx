import { GetConfigs, GetENVConfigs, InitConfig } from 'wailsjs/go/main/CONFIG'
import type { common, main } from 'wailsjs/go/models'

function createDataStore() {
  const [envConfig, setEnvConfig] = createSignal<common.YamlInfo>()
  const [configs, setConfigs] = createStore<main.TypeConfig[]>([])
  async function getEnv() {
    GetENVConfigs().then((res) => {
      setEnvConfig(res)
    })
  }
  async function getData() {
    GetConfigs().then((result) => {
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
  async function refreshData() {
    try {
      await InitConfig()
      await getData()
      // toast.success('加载完成')
    }
    catch (e) {
      // toast.error('加载出错')
      // navigate('/error', {
      //   state: {
      //     msg: e,
      //   },
      // })
    }
  }
  return { configs, getData, refreshData, getEnv, envConfig }
}

export default createRoot(createDataStore)
