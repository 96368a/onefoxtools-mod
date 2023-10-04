import DataStore from '~/store/data'

export default function () {
  const { getEnv, envConfig } = DataStore
  const envs = createMemo<Array<Array<string>>>(() => {
    if (envConfig() && envConfig()!.env)
      return Object.entries(envConfig()!.env)

    return []
  })
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
                    <input class="w-60 bg-base-200 py-2 input rounded-box" value={c[0]}/>
                    <div class="py-2">:</div>
                    <input class="w-full bg-base-200 py-2 input rounded-box" value={c[1]}/>
                    {/* <div class="py-2 bg-base-200 rounded-box w-full">{c[1]}</div> */}
                </div>
                )}</For>
            </div>
        </div>
  )
}
