using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class uniq_email_index_removed : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "IX_mobile_users_email",
                schema: "mobile_api",
                table: "mobile_users");

            migrationBuilder.CreateIndex(
                name: "IX_mobile_users_email",
                schema: "mobile_api",
                table: "mobile_users",
                column: "email");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "IX_mobile_users_email",
                schema: "mobile_api",
                table: "mobile_users");

            migrationBuilder.CreateIndex(
                name: "IX_mobile_users_email",
                schema: "mobile_api",
                table: "mobile_users",
                column: "email",
                unique: true);
        }
    }
}
