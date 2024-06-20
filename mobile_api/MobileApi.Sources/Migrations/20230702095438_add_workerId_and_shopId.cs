using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace MobileApi.Sources.Migrations
{
    /// <inheritdoc />
    public partial class add_workerId_and_shopId : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<DateTime>(
                name: "feedback_date",
                schema: "mobile_api",
                table: "user_feedback",
                type: "timestamp with time zone",
                nullable: false,
                defaultValue: new DateTime(1, 1, 1, 0, 0, 0, 0, DateTimeKind.Unspecified));

            migrationBuilder.AddColumn<int>(
                name: "shop_id",
                schema: "mobile_api",
                table: "user_feedback",
                type: "integer",
                nullable: false,
                defaultValue: 0);

            migrationBuilder.AddColumn<int>(
                name: "worker_id",
                schema: "mobile_api",
                table: "user_feedback",
                type: "integer",
                nullable: false,
                defaultValue: 0);

            migrationBuilder.CreateIndex(
                name: "IX_user_feedback_feedback_date",
                schema: "mobile_api",
                table: "user_feedback",
                column: "feedback_date");

            migrationBuilder.CreateIndex(
                name: "IX_user_feedback_score_quality",
                schema: "mobile_api",
                table: "user_feedback",
                column: "score_quality");

            migrationBuilder.CreateIndex(
                name: "IX_user_feedback_score_service",
                schema: "mobile_api",
                table: "user_feedback",
                column: "score_service");

            migrationBuilder.CreateIndex(
                name: "IX_user_feedback_shop_id",
                schema: "mobile_api",
                table: "user_feedback",
                column: "shop_id");

            migrationBuilder.CreateIndex(
                name: "IX_user_feedback_worker_id",
                schema: "mobile_api",
                table: "user_feedback",
                column: "worker_id");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropIndex(
                name: "IX_user_feedback_feedback_date",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropIndex(
                name: "IX_user_feedback_score_quality",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropIndex(
                name: "IX_user_feedback_score_service",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropIndex(
                name: "IX_user_feedback_shop_id",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropIndex(
                name: "IX_user_feedback_worker_id",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropColumn(
                name: "feedback_date",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropColumn(
                name: "shop_id",
                schema: "mobile_api",
                table: "user_feedback");

            migrationBuilder.DropColumn(
                name: "worker_id",
                schema: "mobile_api",
                table: "user_feedback");
        }
    }
}
