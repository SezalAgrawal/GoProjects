class CreateAccessToken < ActiveRecord::Migration[7.0]
  def change
    create_table :access_tokens, id: :string  do |t|
      t.references :user, type: :string, foreign_key: {to_table: :users}, null: false, index: true
      t.string :token, null: false, unique: true, index: true
      t.boolean :deleted, default: false, null: false
      t.timestamptz :deleted_at
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end
  end
end
