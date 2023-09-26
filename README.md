<p align="center">Wails 的模版，开箱即用的 Vite, Solid, TypeScript 并支持热重载</p>

<p align="center">前端模板来自<a href="https://github.com/96368a/vitesse-lite-solidjs">vitesse-lite-solidjs</a></p>

## 使用模版

```bash
wails init -n my-wails-solid -t https://github.com/96368a/wails-solid-vitesse-template
```

**如果你想使用GoLand调试项目**

```bash
wails init -n my-wails-solid -t https://github.com/96368a/wails-solid-vitesse-template -ide goland
```

## 启动调试

在工程目录中执行 `wails dev` 即可启动。

如果你想在浏览器中调试，请在另一个终端进入 `frontend` 目录，然后执行 `pnpm dev` ，前端开发服务器将在 http://localhost:3333 上运行。

**项目前端默认使用pnpm作为包管理器，如需修改请编辑`wails.json`**

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "solid-vitesse",
  "outputfilename": "solid-vitesse",
  "frontend:install": "pnpm install", 	//更改为npm或者yarn
  "frontend:build": "pnpm build",		//这行也要同步更改
  "frontend:dev:watcher": "pnpm dev",	//这行也要同步更改
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "xxxxxx",
    "email": "xxxxxx@gmail.com"
  }
}
```



## 构建

给你的项目打包，请执行命令： `wails build` 。
