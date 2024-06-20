using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class coinstransactionsadded : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "calc_coins_transaction",
                schema: "mobile_api",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    calc_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_calc_coins_transaction", x => x.id);
                });

            migrationBuilder.CreateTable(
                name: "coins_transaction",
                schema: "mobile_api",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    user_id = table.Column<Guid>(type: "uuid", nullable: false),
                    date = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    transaction_type = table.Column<string>(type: "text", nullable: false),
                    zebra_coins = table.Column<decimal>(type: "numeric", nullable: false),
                    note = table.Column<string>(type: "text", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_coins_transaction", x => x.id);
                });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "calc_coins_transaction",
                schema: "mobile_api");

            migrationBuilder.DropTable(
                name: "coins_transaction",
                schema: "mobile_api");
        }
    }
}
