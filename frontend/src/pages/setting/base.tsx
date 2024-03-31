import toast from 'solid-toast'
import { WindowSetTitle } from 'wailsjs/runtime/runtime'
import DataStore from '~/store/data'

export default function () {
  const { envConfig, saveEnv, refreshEnv, updateEnvConfig } = DataStore

  const [title, setTitle] = createSignal(envConfig.title)
  function SaveENVConfigs() {
    toast.promise(
      (async function () {
        await updateEnvConfig('title', title())
        await saveEnv()
        WindowSetTitle(title())
      })(),
      {
        loading: '保存中...',
        success: () => <span> 保存成功</span>,
        error: '保存出错',
      },
    )
  }
  return (
        <div class="h-full w-full bg-base-300 px-4 py-4 rounded-box">
            <div class="py-2">
                基础设置
            </div>
            <div class="flex flex-col gap-2 px-2">
                <div class="flex gap-4">
                    <input class="w-60 bg-base-200 py-2 input input-bordered rounded-box" readonly value="应用标题" />
                    <div class="py-2">:</div>
                    <input type="text" class="flex-1 bg-base-200 py-2 input rounded-box" value={envConfig.title} onchange={e => setTitle(e.currentTarget.value)}/>
                </div>
            </div>
            <div class='flex justify-center gap-4 py-4'>

                <button class='px-10 btn' onclick={SaveENVConfigs}>保存配置</button>
                <button class='px-10 btn btn-warning' onclick={refreshEnv}>重载配置</button>
            </div>
        </div>
  )
}
