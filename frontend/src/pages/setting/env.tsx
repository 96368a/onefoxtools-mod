import toast from 'solid-toast'
import DataStore from '~/store/data'

export default function () {
  const { getEnv, envConfig, saveEnv, refreshEnv, updateEnv } = DataStore
  const envs = createMemo(() => {
    if (envConfig && envConfig!.env)
      return Object.entries(envConfig!.env)

    return []
  })
  function SaveENVConfigs() {
    toast.promise(
      saveEnv(),
      {
        loading: '保存中...',
        success: () => <span> 保存成功</span>,
        error: '保存出错',
      },
    )
  }
  function RefreshEnv() {
    toast.promise(
      (async () => {
        refreshEnv()
      })(),
      {
        loading: '加载中...',
        success: () => <span> 加载成功</span>,
        error: '加载出错',
      },
    )
  }
  onMount(async () => {
    await getEnv()
  })
  return (
    <div class="h-full w-full bg-base-300 px-4 py-4 rounded-box">
      <div class="py-2" onclick={getEnv}>
        环境变量设置
      </div>
      <div class="flex flex-col gap-2 px-2">
        <For each={envs()}>{c => (
          <div class="flex gap-4">
            <input class="w-60 bg-base-200 py-2 input rounded-box" readonly value={c[0]} />
            <div class="py-2">:</div>
            <select class="w-full bg-base-200 input select select-bordered rounded-box" onchange={e => updateEnv(c[0], Number.parseInt(e.currentTarget.value))}>
              <For each={c[1].list}>
                {(v, i) => <option value={i()}>{v}</option>}
              </For>
            </select>
            {/* <input class="w-full bg-base-200 py-2 input rounded-box" value={c[1]} onchange={e => updateEnv(c[0], e.currentTarget.value)}/> */}
            {/* <div class="py-2 bg-base-200 rounded-box w-full">{c[1]}</div> */}
          </div>
        )}</For>
      </div>
      <div class='flex justify-center gap-4 py-4'>

        <button class='px-10 btn' onclick={SaveENVConfigs}>保存配置</button>
        <button class='px-10 btn btn-warning' onclick={RefreshEnv}>重载配置</button>
      </div>
    </div>
  )
}
