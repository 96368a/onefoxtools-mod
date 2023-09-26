/// <reference types="vitest" />

import path from 'node:path'
import { defineConfig } from 'vite'
import solidPlugin from 'vite-plugin-solid'
import Pages from 'vite-plugin-pages'
import Unocss from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'

export default defineConfig({
  resolve: {
    alias: {
      '~/': `${path.resolve(__dirname, 'src')}/`,
      '@wailsjs/': `${path.resolve(__dirname, 'wailsjs')}/`,
    },
  },
  plugins: [
    solidPlugin(),
    Pages(),
    Unocss(),
    AutoImport({
      imports: [
        'solid-js',
        '@solidjs/router',
      ],
      dts: true,
      dirs: [
        'src/primitives',
      ],
    }),

  ],
  build: {
    target: 'esnext',
  },
  test: {
    environment: 'jsdom',
    transformMode: {
      web: [/.[jt]sx?/],
    },
    threads: false,
    isolate: false,
  },
})
