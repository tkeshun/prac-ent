entのチュートリアルやる
よくわかんないところは補完してく

Quick Introduction

データベーススキーマをグラフ構造として簡単にモデリング  
スキーマをプログラム的なGoコードとして定義  
コード生成に基づく静的型付け  
データベースクエリやグラフのトラバーサルが簡単に記述可能  
Goテンプレートを使用して簡単に拡張やカスタマイズが可能

- プロジェクト初期化

    go mod init tutorial

- インストール  
    プロジェクトへの導入  
    go get entgo.io/ent/cmd/ent

    EntCLIを使えるようにするため、グローバルにもインストールする
    go install entgo.io/ent/cmd/ent@latest


- Create Your First Schema

    go run -mod=mod entgo.io/ent/cmd/ent new User

    EntCLIを使って新しいエンティティ（User）のスキーマを生成するコード

    コマンド実行後にent/が生成される
    ```
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

  - Fields(): データベーステーブルのカラムに対応するフィールドを定義
  - Edges(): 他のエンティティとのリレーションを定義(1対多、多対多など) 


Userスキーマに２つフィールドを追加する
```
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

CRUD関連のコード生成

```
go generate ./ent
```

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

pgxを使うときはコネクションの開き方が変わる
ent.Openだと
```
go run main.go 
2024/11/20 01:16:18 failed opening connection to postgres: sql: unknown driver "postgres" (forgotten import?)
exit status 1
```
となり接続できない

解決策：https://entgo.io/ja/docs/sql-integration/

以下サンプルのようにsql.Openでコネクションを開いてからEntに渡す。

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

- 関数化

```
// Userを１件作って、作ったエンティティをリターンする
func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
    u, err := client.User.
        Create().
        SetAge(30).
        SetName("a8m").
        Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed creating user: %w", err)
    }
    log.Println("user was created: ", u)
    return u, nil
}
```

- Relationの定義

    新たにCarテーブルを定義する
    テンプレート生成
    `go run -mod=mod entgo.io/ent/cmd/ent new Car`
    ent/schema/car.goに追記
    フィールド（カラムを定義する）
    ※ PostgresqlにText型で保存する方法：.SchemaType(map[string]string{"postgres": "text"})をつける。デフォルトではVARCHAR型なので、変更する必要がある。
    
    ```
    func (Car) Fields() []ent.Field {
        return []ent.Field{
            field.String("model").SchemaType(map[string]string{"postgres": "text"}),
            field.String("car_number").SchemaType(map[string]string{"postgres": "text"})
            field.Time("registered_at").Default(time.Now), // Defaultでデフォルト値を指定できる
        }
    }
    ```

    その他の制約系メソッド
    - 必須長の指定
        .NotEmpty()
    - 正規表現によるバリデーション
        .Match(regexp.MustCompile("[a-zA-Z_]+$")
    - 最大長指定
        .MaxLen(50)

    schema/user.goにCarとの関係を記述する
    ```
    // Edges of the User.
    func (User) Edges() []ent.Edge {
        return []ent.Edge{
            edge.To("cars", Car.Type),
        }
    }
    ```
    再生成
    
    `go generate ./ent`

    - 外部キーの名称を変更したい場合
    参考：https://mackerel.io/ja/blog/entry/2024/01/10/162149

    指定しないとテーブル名を結合したものになる
    参照元と参照先のそれぞれで設定することで名前を変えられる。
    デフォの外部キーはNullableみたいなので、設定必須かも
    あとから変えてももとのカラムは消えないみたいなので、最初にやるとよい

    