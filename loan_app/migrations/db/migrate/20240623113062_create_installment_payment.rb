class CreateInstallmentPayment < ActiveRecord::Migration[7.0]
  def change
    create_table :installment_payments, id: :string  do |t|
      t.references :loan, type: :string, foreign_key: {to_table: :loans}, null: false, index: true
      t.references :installment, type: :string, foreign_key: {to_table: :installments}, null: false, index: true
      t.integer :amount_in_cents, null: false
      t.string :one_time_settlement_id, null: false
      t.timestamptz :created_at, null: false
      t.timestamptz :updated_at, null: false
    end
  end
end
