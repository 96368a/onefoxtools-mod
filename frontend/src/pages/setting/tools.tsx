import { GenerateConfig } from 'wailsjs/go/main/GOContext'

export default () => {
  function genConfig() {
    GenerateConfig().catch((err) => {
      console.error(err)
    })
  }
  return (
        <div class="h-full w-full bg-base-300 px-4 py-4 rounded-box">
            <div class="py-2">
                工具设置
            </div>
            <div class='flex justify-center gap-4 py-4'>
                <button class='px-10 btn btn-warning' onclick={genConfig}>初始化工具配置</button>
            </div>
        </div>
  )
}
