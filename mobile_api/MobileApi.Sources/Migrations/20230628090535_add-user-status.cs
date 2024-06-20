using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class adduserstatus : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<DateTime>(
                name: "removed_date",
                schema: "mobile_api",
                table: "mobile_users",
                type: "timestamp with time zone",
                nullable: true);

            migrationBuilder.AddColumn<int>(
                name: "status",
                schema: "mobile_api",
                table: "mobile_users",
                type: "integer",
                nullable: false,
                defaultValue: 0);

            migrationBuilder.CreateIndex(
                name: "IX_mobile_users_status",
                schema: "mobile_api",
                table: "mobile_users",
                column: "status");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "IX_mobile_users_status",
                schema: "mobile_api",
                table: "mobile_users");

            migrationBuilder.DropColumn(
                name: "removed_date",
                schema: "mobile_api",
                table: "mobile_users");

            migrationBuilder.DropColumn(
                name: "status",
                schema: "mobile_api",
                table: "mobile_users");
        }
    }
}
