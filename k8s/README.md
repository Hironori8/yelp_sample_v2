# Yelp Sample Kubernetes Deployment

Docker ComposeからKubernetesに移行したYelp Sampleアプリケーションのデプロイメント設定です。

## ディレクトリ構成

```
k8s/
├── base/                    # ベースマニフェスト
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── pvc.yaml
│   ├── postgres-deployment.yaml
│   ├── postgres-service.yaml
│   ├── cassandra-deployment.yaml
│   ├── cassandra-service.yaml
│   ├── gateway-deployment.yaml
│   ├── gateway-service.yaml
│   ├── *-service-deployment.yaml
│   ├── microservices-services.yaml
│   └── kustomization.yaml
├── overlays/
│   ├── dev/                 # 開発環境用設定
│   └── prod/                # 本番環境用設定
└── README.md
```

## デプロイ方法

### 前提条件
- Kubernetesクラスターが準備済み
- `kubectl`がインストール済み
- `kustomize`がインストール済み（またはkubectl 1.14+）

### デプロイ手順

#### 開発環境
```bash
# 開発環境にデプロイ
kubectl apply -k k8s/overlays/dev/

# デプロイ状況確認
kubectl get pods -n yelp-app

# サービス確認
kubectl get services -n yelp-app
```

#### 本番環境
```bash
# 本番環境にデプロイ
kubectl apply -k k8s/overlays/prod/

# デプロイ状況確認
kubectl get pods -n yelp-app

# サービス確認
kubectl get services -n yelp-app
```

### イメージビルドとプッシュ

```bash
# Dockerイメージをビルド
docker build -f services/gateway/Dockerfile -t your-registry/yelp-gateway:v1.0.0 .
docker build -f services/business/Dockerfile -t your-registry/yelp-business-service:v1.0.0 .
docker build -f services/review/Dockerfile -t your-registry/yelp-review-service:v1.0.0 .
docker build -f services/logging/Dockerfile -t your-registry/yelp-logging-service:v1.0.0 .
docker build -f services/auth/Dockerfile -t your-registry/yelp-auth-service:v1.0.0 .

# レジストリにプッシュ
docker push your-registry/yelp-gateway:v1.0.0
docker push your-registry/yelp-business-service:v1.0.0
docker push your-registry/yelp-review-service:v1.0.0
docker push your-registry/yelp-logging-service:v1.0.0
docker push your-registry/yelp-auth-service:v1.0.0
```

## 設定のカスタマイズ

### 環境変数の変更
`k8s/base/configmap.yaml`と`k8s/base/secret.yaml`を編集してください。

### リソース制限の変更
各環境の`resource-patch.yaml`を編集してください。

### レプリカ数の変更
各環境の`replica-patch.yaml`を編集してください。

## トラブルシューティング

### ログ確認
```bash
# Pod の ログ確認
kubectl logs -f deployment/gateway -n yelp-app

# 特定のPodのログ確認
kubectl logs <pod-name> -n yelp-app
```

### デバッグ
```bash
# Pod内に接続
kubectl exec -it <pod-name> -n yelp-app -- /bin/bash

# サービスの詳細確認
kubectl describe service gateway-service -n yelp-app
```

### クリーンアップ
```bash
# 開発環境削除
kubectl delete -k k8s/overlays/dev/

# 本番環境削除
kubectl delete -k k8s/overlays/prod/

# Namespace全体削除
kubectl delete namespace yelp-app
```

## 主な特徴

- **Kustomize**: 環境別設定の管理
- **リソース制限**: CPU/メモリ使用量の制御
- **ヘルスチェック**: Liveness/Readiness Probe
- **永続化**: PostgreSQL/Cassandra用PVC
- **LoadBalancer**: Gateway用外部アクセス
- **ConfigMap/Secret**: 設定とシークレットの分離