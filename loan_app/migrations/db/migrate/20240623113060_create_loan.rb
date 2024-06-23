class CreateLoan < ActiveRecord::Migration[7.0]
  def change
    create_table :loans, id: :string  do |t|
      t.references :user, type: :string, foreign_key: {to_table: :users}, null: false, index: true
      t.integer :amount_in_cents, null: false
      t.integer :term, null: false
      t.integer :frequency_in_days, null: false
      t.string :status, null: false
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end
  end
end
