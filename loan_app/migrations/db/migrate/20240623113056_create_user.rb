class CreateUser < ActiveRecord::Migration[7.0]
  def change
    create_table :users, id: :string  do |t|
      t.string :name, null: false
      t.string :hashed_password, null: false
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end

    add_index(:users, [:name, :hashed_password], unique: true)
  end
end
