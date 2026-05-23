/**
 * @vitest-environment jsdom
 */
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import AppFooter from "components/layout/AppFooter";

describe("AppFooter", () => {
  it("renders external footer links safely", () => {
    render(
      <AppFooter
        githubUrl="https://github.com/example"
        linkedinUrl="https://linkedin.com/in/example"
        websiteUrl="https://example.com"
      />,
    );

    expect(screen.getByRole("link", { name: "GitHub" })).toHaveAttribute(
      "href",
      "https://github.com/example",
    );
    expect(screen.getByRole("link", { name: "LinkedIn" })).toHaveAttribute(
      "rel",
      "noopener noreferrer",
    );
    expect(screen.getByRole("link", { name: "Website" })).toHaveAttribute(
      "target",
      "_blank",
    );
  });
});
