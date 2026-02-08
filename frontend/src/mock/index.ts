import type { MockMethod } from 'vite-plugin-mock'

import albums from './albums'
import anniversaries from './anniversaries'
import dashboard from './dashboard'
import moments from './moments'
import photos from './photos'
import places from './places'
import systemInit from './system'
import upload from './upload'
import users from './users'
import wishes from './wishes'

export { default as albums } from './albums'
export { default as anniversaries } from './anniversaries'
export { default as dashboard } from './dashboard'
export { default as moments } from './moments'
export { default as places } from './places'
export { default as system } from './system'
export { default as upload } from './upload'
export { default as users } from './users'
export { default as wishes } from './wishes'

export default [
  ...systemInit,
  ...anniversaries,
  ...places,
  ...moments,
  ...albums,
  ...photos,
  ...wishes,
  ...users,
  ...dashboard,
  ...upload,
] as MockMethod[]
