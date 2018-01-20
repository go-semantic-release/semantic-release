const os = require('os')
const pkg = require('./package.json')

const rootUrl = 'https://get-release.xyz/go-semantic-release/semantic-release/'
function getPlatformArch (a, p) {
  const platform = {
    win32: 'windows'
  }
  const arch = {
    x64: 'amd64',
    x32: '386'
  }
  return (platform[p] ? platform[p] : p) + '/' + (arch[a] ? arch[a] : a)
}

console.log('downloading file....')
require('download')(rootUrl + getPlatformArch(os.arch(), os.platform()) + '/' + pkg.version)
.pipe(require('fs').createWriteStream(require('path').join(__dirname, 'bin'), {mode: 0o755}))
