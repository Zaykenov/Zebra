using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class feedbacks_added : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "user_feedback",
                schema: "mobile_api",
                columns: table => new
                {
                    Id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    user_id = table.Column<Guid>(type: "uuid", nullable: false),
                    check_id = table.Column<int>(type: "integer", nullable: false),
                    score_quality = table.Column<double>(type: "double precision", nullable: false),
                    score_service = table.Column<double>(type: "double precision", nullable: false),
                    feedback_text = table.Column<string>(type: "text", nullable: true),
                    check_json = table.Column<string>(type: "text", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_user_feedback", x => x.Id);
                });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "user_feedback",
                schema: "mobile_api");
        }
    }
}
