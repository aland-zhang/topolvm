PoC
====

### 目的
CSIプラグインの既存実装がnecoチームで必要な機能・水準を満たしていないと判断したため、
CSIプラグインをnecoで実装することを検討する。しかし、CSI実装のための知識が不十分なため、
本PoCで必要な知識の整理を行う。

### 作業内容

#### 既存実装の動作確認
PVC/PV/snapshotを実装している[既存のCSIプラグイン](https://github.com/kubernetes-csi/csi-driver-host-path/tree/master/examples)を`neco/dctest`環境上で動作確認した。

- CSIプラグインをデプロイするのに必要な各サイドカーコンテナ(`external-provisioner/external-attacher`/`external-snapshotter`)の[RBACマニフェスト](https://github.com/kubernetes-csi/csi-driver-host-path/blob/670a96b163a3d605ef8520c61be743a1aa05894c/README.md#deployment)を取得。
- hostpathプラグインをデプロイする[マニフェスト](https://github.com/kubernetes-csi/csi-driver-host-path/tree/master/deploy/kubernetes-1.14/hostpath)を取得。
  - `hostpathplugin`コンテナは`topolvm`リポジトリを自前でビルドしたものを`quay.io/cybozu/csi-hostpath:latest`に用意した。
```console
make hostpath
docker build -t quay.io/cybozu/csi-hostpath:latest .
docker login quay.io
docker push quay.io/cybozu/csi-hostpath:latest
```
　- [対応するコンテナイメージ](https://github.com/kubernetes-csi/csi-driver-host-path/blob/670a96b163a3d605ef8520c61be743a1aa05894c/deploy/kubernetes-1.14/hostpath/csi-hostpath-plugin.yaml#L66)の箇所を差し替えた。
  - `kubelet`のルートファイルシステムは原則Readonlyであるため、[PVを作成するディレクトリ](https://github.com/kubernetes-csi/csi-driver-host-path/blob/670a96b163a3d605ef8520c61be743a1aa05894c/deploy/kubernetes-1.14/hostpath/csi-hostpath-plugin.yaml#L135)を`/tmp/csi-hostpath-data`に変更した。
- 以上を反映した各マニフェストを`deploy`上に配置し、また動作確認用のマニフェストを`example`上に配置した。

#### 再現方法
- `hostpathplugin`が`previllege`権限を要求しており、necogcpのデフォルトの`PodSecurityPolicy`では権限がなく動作不可であった。`PodSecurityPolicy`を空配列にし、無効にした(ブートサーバ上で`ckecli cluster get > manifest.yaml`とし、生成されたマニフェストを変更後に`ckecli cluster set manifest.yaml`)。

- gcp環境上のブートサーバーに`deploy`および`example`を`dcscp`等で配置し、各マニフェストをapplyした。
```console
kubectl apply -f deploy/rbac
kubectl apply -f deploy
kubectl apply -f example
```
- 上記によりDynamic Provisioningが行われていることを確認した。snapshotは既存実装が簡易的なものであったため動作しなかった。
