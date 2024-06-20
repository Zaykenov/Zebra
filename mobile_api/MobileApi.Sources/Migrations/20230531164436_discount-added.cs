using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class discountadded : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<decimal>(
                name: "discount",
                schema: "mobile_api",
                table: "mobile_users",
                type: "numeric",
                nullable: false,
                defaultValue: 0m);

            migrationBuilder.AddColumn<decimal>(
                name: "zebra_coin_balance",
                schema: "mobile_api",
                table: "mobile_users",
                type: "numeric",
                nullable: false,
                defaultValue: 0m);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "discount",
                schema: "mobile_api",
                table: "mobile_users");

            migrationBuilder.DropColumn(
                name: "zebra_coin_balance",
                schema: "mobile_api",
                table: "mobile_users");
        }
    }
}
