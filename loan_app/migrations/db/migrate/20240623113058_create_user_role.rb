class CreateUserRole < ActiveRecord::Migration[7.0]
  def change
    create_table :user_roles, id: :string  do |t|
      t.references :user, type: :string, foreign_key: {to_table: :users}, null: false, index: true
      t.references :role, type: :string, foreign_key: {to_table: :roles}, null: false, index: false
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end
  end
end
