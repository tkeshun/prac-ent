調べて埋める

## Entの概要と背景

### Entとは何か

https://github.com/ent/ent

- Schema As Code  
任意のデータベーススキーマをGoオブジェクトとしてモデリングできる

- グラフ構造の容易なトラバース  
クエリや集計を実行し、どんなグラフ構造でも簡単に辿れる。

- 静的型付けされた明示的なAPI  
コード生成を活用し、100%静的型付けされた明示的なAPIを提供する

- 複数ストレージドライバのサポート  
MySQL、MariaDB、TiDB、PostgreSQL、CockroachDB、SQLite、Gremlinに対応している

- 拡張性  
Goテンプレートを使用して簡単に拡張やカスタマイズが可能

### Entの使用が推奨されるシナリオ

## Entのセットアップ

- プロジェクト初期化

`go mod init tutorial`

- インストール

プロジェクトへの導入

`go get entgo.io/ent/cmd/ent`

EntCLIを使えるようにするため、グローバルにもインストールする
`go install entgo.io/ent/cmd/ent@latest`


### 初期設定と簡単なモデル作成

テンプレートの作成

`go run -mod=mod entgo.io/ent/cmd/ent new <テーブル名>`

例
```
go run -mod=mod entgo.io/ent/cmd/ent new User

EntCLIを使って新しいエンティティ（User）のスキーマを生成するコード

コマンド実行後にent/が生成される

├── do-ent-tutorial.md
├── ent
│   ├── generate.go
│   └── schema
│       └── user.go
├── go.mod
└── go.sum
```


- ent/generate.go

Ent CLIによるコード生成に必要なエントリポイント
このファイルは、コード生成時にプロジェクト内のスキーマを指定して実行される

- ent/schema/user.go

Userエンティティのスキーマファイル
エンティティのフィールド（列）やリレーション（エッジ）を定義する
エッジとはデータベースでいうところの「テーブル同士のリレーション（外部キー）」に相当する

```
Fields(): データベーステーブルのカラムに対応するフィールドを定義
Edges(): 他のエンティティとのリレーションを定義(1対多、多対多など)
```


## 基本的なモデル定義

EntCLIで生成した基本のテンプレートに記述する形でカラム名やテーブル間の関係を決める


- Fields(): データベーステーブルのカラムに対応するフィールドを定義
- Edges(): 他のエンティティとのリレーションを定義(1対多、多対多など)

```
例）
Userスキーマに２つフィールドを追加する

// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Positive(),
        field.String("name").
            Default("unknown"),
    }
}
```

- CRUD関連のコード生成

    `go generate ./ent`

実行後のディレクトリ構造

```
.
├── do-ent-tutorial.md
├── ent
│   ├── client.go
│   ├── ent.go // Entのメタデータや設定を保持するファイル
│   ├── enttest
│   │   └── enttest.go // テストヘルパー
│   ├── generate.go // コード生成のためのエンドポイント
│   ├── hook
│   │   └── hook.go // CRUD操作やクエリの前後にカスタムロジックを挿入する場合に使用
│   ├── migrate
│   │   ├── migrate.go // スキーマのマイグレーションを実行するためのロジック
│   │   └── schema.go // スキーマに関するメタデータやマイグレーション用の設定
│   ├── mutation.go // CRUD操作の際に使用される内部的な構造体を定義
│   ├── predicate
│   │   └── predicate.go // クエリにおけるフィルタリング条件（Predicate）を定義
│   ├── runtime
│   │   └── runtime.go // Entの実行時に必要なメタデータや初期化処理
│   ├── runtime.go
│   ├── schema
│   │   └── user.go // User（定義したやつ）
│   ├── tx.go // トランザクション（Tx）の管理機能
│   ├── user
│   │   ├── user.go // ユーザーのCRUD操作に必要な構造体やメソッド
│   │   └── where.go // フィルタリング条件（WHERE句）を簡単に記述するためのメソッド
│   ├── user.go
│   ├── user_create.go
│   ├── user_delete.go
│   ├── user_query.go
│   └── user_update.go
├── go.mod
└── go.sum
```


### フィールドタイプと制約


### リレーション (One-to-One, One-to-Many, Many-to-Many)


### プリミティブフィールドとカスタムフィールドの作成

## DB接続 （pgx）

sql.Openでコネクションを開いてからEntに渡す

```
package main

import (
    "context"
    "database/sql"
    "log"

    "<project>/ent"

    "entgo.io/ent/dialect"
    entsql "entgo.io/ent/dialect/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

// Open new connection
func Open(databaseUrl string) *ent.Client {
    db, err := sql.Open("pgx", databaseUrl)
    if err != nil {
        log.Fatal(err)
    }

    // Create an ent.Driver from `db`.
    drv := entsql.OpenDB(dialect.Postgres, db)
    return ent.NewClient(ent.Driver(drv))
}

func main() {
    client := Open("postgresql://user:password@127.0.0.1/database")

    // Your code. For example:
    ctx := context.Background()
    if err := client.Schema.Create(ctx); err != nil {
        log.Fatal(err)
    }
    users, err := client.User.Query().All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(users)
}
```

## シンプルなCRUD操作


### エンティティの作成、読み取り、更新、削除


### クエリビルダーの基本的な使い方

クエリビルダーは使い回さない
クエリは実行後も保持されるため、条件を変えて複数回使うときはイテレーションごとに定義しなおす
参考：　https://qiita.com/k_t_rat/items/1b7344b1d7531d190d1b



### エラー処理とロギングのベストプラクティス


## 高度なクエリ構築
動的フィルタリング
集計 (Aggregation) とグループ化
SQL関数やカスタムクエリの統合

## 2.2 リレーションシップの管理
外部キーとリレーションの扱い
Eager Loading vs Lazy Loading
リレーションの再帰的クエリ

## 2.3 ミドルウェアとフック
ミドルウェアの仕組みと実装
フックの活用: BeforeCreate、AfterUpdateなど
トランザクション管理

## 2.4 スキーママイグレーション
Entの自動マイグレーション
カスタムマイグレーションの作成
本番環境での安全なマイグレーション手法


## Atlasとの連携

1.1 Atlasの概要と特徴
Atlasとは何か
スキーマ管理ツールとしての位置付け
他のマイグレーションツール（Flyway, Liquibaseなど）との比較
1.2 Atlasのセットアップ
必要な環境の準備（Docker, Goなど）
Atlas CLIのインストール
簡単なスキーマ作成と適用
1.3 スキーマ管理の基本
スキーマファイルの記述方法（HCL形式）
スキーマの初期化とバージョン管理
マイグレーションの作成と適用
1.4 スキーマの差分とデプロイ
既存データベースとの差分検出
ローカル環境でのシミュレーション
本番環境への安全なデプロイ方法

2. EntとAtlasの統合
2.1 EntとAtlasの役割分担
Entによるモデル設計
Atlasによるスキーマ管理とマイグレーション
双方を組み合わせるメリット
2.2 EntスキーマからAtlasスキーマへの変換
EntコードからAtlas HCLファイルを生成
カスタムフィールドや制約の反映
マイグレーション自動生成の手順
2.3 双方向の同期と管理
EntとAtlasの間でのスキーマ差分管理
スキーマ変更が多いプロジェクトでの運用
リアルタイムにスキーマを同期するベストプラクティス
2.4 EntマイグレーションとAtlas CLIの統合
EntのmigrateパッケージとAtlas CLIの使い分け
CI/CDパイプラインでの統合手法
ロールバック戦略の構築

3.2 カスタムバリデーションとガバナンス
独自のスキーマポリシーをAtlasで設定
フィールド名やデータ型の一貫性を確保する方法
チーム規約に基づいたスキーマの自動検証

## 2.5 テスト
モックを使ったユニットテスト
インメモリデータベースを活用したテスト戦略
実際のデータベースを使用した統合テスト

## 3.1 パフォーマンスチューニング
クエリの最適化
インデックスの活用
キャッシュ戦略と分散キャッシュの統合

## 3.2 複雑なユースケース
マルチテナンシー対応
JSON型や他の特殊なデータ型の使用
地理空間クエリの実装

## 3.3 エコシステムとの統合
OpenTelemetryを使ったトレース
GraphQLサーバーの自動生成
REST APIとの統合

## 3.4 本番環境の運用
本番データベースの監視とメンテナンス
ロールバック戦略
CI/CDパイプラインへのEntの統合

4. EntとAtlasを組み合わせたプロジェクト演習
4.1 小規模プロジェクト
単一データベースのCRUDアプリケーションを構築
EntモデルからAtlasスキーマを生成し、運用までを体験
4.2 大規模プロジェクト
複雑なリレーションを持つスキーマを管理
EntのカスタムロジックとAtlasのスキーマポリシーを活用
複数の開発者間でのスキーマ変更を調整
4.3 実運用シミュレーション
スキーマの大幅変更を伴うシナリオでの運用
スキーマ変更のレビューとCI/CD統合
本番データベースでのAtlasの活用