using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class usersimpler : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.EnsureSchema(
                name: "mobile_api");

            migrationBuilder.CreateTable(
                name: "mobile_users",
                schema: "mobile_api",
                columns: table => new
                {
                    id = table.Column<Guid>(type: "uuid", nullable: false),
                    email = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: false),
                    name = table.Column<string>(type: "text", nullable: false),
                    birth_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: true),
                    reg_date = table.Column<DateTime>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_mobile_users", x => x.id);
                });

            migrationBuilder.CreateIndex(
                name: "IX_mobile_users_email",
                schema: "mobile_api",
                table: "mobile_users",
                column: "email",
                unique: true);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "mobile_users",
                schema: "mobile_api");
        }
    }
}
