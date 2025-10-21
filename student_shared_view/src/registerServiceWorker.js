/* eslint-disable no-console */

import { register } from 'register-service-worker'

// 只在生产环境且支持Service Worker的浏览器中注册
if (process.env.NODE_ENV === 'production' && 'serviceWorker' in navigator) {
  register(`${process.env.BASE_URL}service-worker.js`, {
    ready () {
      console.log(
        'App is being served from cache by a service worker.\n' +
        'For more details, visit https://goo.gl/AFskqB'
      )
    },
    registered (registration) {
      console.log('Service worker has been registered.', registration)
    },
    cached () {
      console.log('Content has been cached for offline use.')
    },
    updatefound () {
      console.log('New content is downloading.')
    },
    updated (registration) {
      console.log('New content is available; please refresh.', registration)
    },
    offline () {
      console.log('No internet connection found. App is running in offline mode.')
    },
    error (error) {
      console.error('Error during service worker registration:', error)
      // 避免在开发环境中显示错误
      if (process.env.NODE_ENV !== 'development') {
        console.error('Service Worker registration failed:', error)
      }
    }
  })
} else if (process.env.NODE_ENV === 'development') {
  // 在开发环境中禁用Service Worker相关功能
  console.log('Service Worker is disabled in development mode')
}
