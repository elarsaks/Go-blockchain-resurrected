/**
 * @vitest-environment jsdom
 */
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Notification from "components/shared/Notification";

describe("Notification", () => {
  it("renders the message and info loader for active info notifications", () => {
    const { container } = render(
      <Notification
        type="info"
        message="Fetching blockchain data..."
        insideContainer={false}
      />,
    );

    expect(screen.getByText("Fetching blockchain data...")).toBeInTheDocument();
    expect(container.querySelector("img")).toBeInTheDocument();
  });

  it("does not render empty notifications", () => {
    const { container } = render(
      <Notification type="success" message="" insideContainer={true} />,
    );

    expect(container).toBeEmptyDOMElement();
  });
});
