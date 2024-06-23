class CreateInstallment < ActiveRecord::Migration[7.0]
  def change
    create_table :installments, id: :string  do |t|
      t.references :loan, type: :string, foreign_key: {to_table: :loans}, null: false, index: true
      t.integer :amount_in_cents, null: false
      t.integer :serial_no, null: false
      t.string :status, null: false
      t.timestamptz :due_date, null: false
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end
  end
end
