class CreateRole < ActiveRecord::Migration[7.0]
  def change
    create_table :roles, id: :string  do |t|
      t.string :name, null: false, unique: true
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end
  end
end
