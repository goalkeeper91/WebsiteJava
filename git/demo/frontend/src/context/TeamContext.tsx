import React, { createContext, useContext, useEffect, useState } from 'react';
import { useAuth } from './AuthContext';
import { getManagedChannels } from '../features/team/api';
import type { ManagedChannel } from '../features/team/types';
import { setActingAsChannel as setApiActingAsChannel, setForbiddenHandler } from '../lib/apiFetch';

type TeamContextType = {
  actingAsChannel: ManagedChannel | null;
  managedChannels: ManagedChannel[];
  setActingAsChannel: (channel: ManagedChannel | null) => void;
  refreshManagedChannels: () => Promise<void>;
};

const TeamContext = createContext<TeamContextType>({
  actingAsChannel: null,
  managedChannels: [],
  setActingAsChannel: () => {},
  refreshManagedChannels: async () => {},
});

export const TeamProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, authChecked } = useAuth();
  const [actingAsChannel, setActingAsChannelState] = useState<ManagedChannel | null>(null);
  const [managedChannels, setManagedChannels] = useState<ManagedChannel[]>([]);

  const refreshManagedChannels = async () => {
    try {
      const channels = await getManagedChannels();
      setManagedChannels(channels);
    } catch (err) {
      console.error('Fehler beim Laden der verwaltbaren Kanäle:', err);
    }
  };

  useEffect(() => {
    if (authChecked && isAuthenticated) {
      refreshManagedChannels();
    }
  }, [authChecked, isAuthenticated]);

  useEffect(() => {
    setApiActingAsChannel(actingAsChannel?.owner_twitch_id ?? null);
  }, [actingAsChannel]);

  useEffect(() => {
    // If the backend rejects a request for the currently-managed channel
    // (e.g. access was just revoked), fall back to the own channel instead
    // of silently keeping every subsequent request 403ing.
    setForbiddenHandler(() => setActingAsChannelState(null));
    return () => setForbiddenHandler(null);
  }, []);

  return (
    <TeamContext.Provider
      value={{
        actingAsChannel,
        managedChannels,
        setActingAsChannel: setActingAsChannelState,
        refreshManagedChannels,
      }}
    >
      {children}
    </TeamContext.Provider>
  );
};

export const useTeam = () => useContext(TeamContext);
