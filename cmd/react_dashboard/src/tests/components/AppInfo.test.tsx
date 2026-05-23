/**
 * @vitest-environment jsdom
 */
import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import AppInfo from "components/layout/AppInfo";

describe("AppInfo", () => {
  it("describes the desktop wallet layout and links to GitHub", () => {
    window.innerWidth = 1024;

    render(<AppInfo />);

    expect(
      screen.getByText(/wallet on the left represents a miner/i),
    ).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "GitHub" })).toHaveAttribute(
      "href",
      "https://github.com/elarsaks/Go-blockchain-resurrected",
    );
  });

  it("updates the layout description for mobile widths", () => {
    window.innerWidth = 700;

    render(<AppInfo />);
    fireEvent(window, new Event("resize"));

    expect(
      screen.getByText(/wallet on the up represents a miner/i),
    ).toBeInTheDocument();
  });
});
